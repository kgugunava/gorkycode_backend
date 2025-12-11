from OSRM.location import Location
from ortools.constraint_solver import routing_enums_pb2
from ortools.constraint_solver import pywrapcp
import numpy as np
import requests
import logging

class RoutePlanner:
    def loadTime(self, from_id, to_id):

        try:
            cursor = self.connection.cursor()
            cursor.execute(f'''
                    SELECT duration_hours
                    FROM walking_matrix
                    WHERE from_id = {from_id}
                    AND to_id = {to_id};
                ''')

            duration = cursor.fetchone()

            # print(duration[0])
            cursor.close()
        except:
            self.logger.warning("TIME EXTRACTION PROBLEM ")
    
        # +1 потому что срезает минуты, а мы должны быть уверены,
        # что он не вылезет за лимит времени
        return int(duration[0] * 60) + 31

    def __init__(self, connection, user_point: Location,
                 T_limit: int, locations, weight_rerank, weight_distance):
        
        self.logger = logging.getLogger("ROUTE_PLANNER")
        self.logger.debug("ROUTE_PLANNER IS LAUNCHED")

        self.connection = connection
        self.time_limit = T_limit  # в минутах
        # тут первая точка - это всегда местоположение пользователя!
        # (потому что это будет старт) ТАК ЖЕ ДЛЯ НЕЁ НАДО
        # ЗАРАНЕЕ ПОСЧИТАТЬ ВРЕМЯ ДЛЯ ВСЕХ ВЫБРАННЫХ ТОЧЕК!!!
        self.all_nodes = [user_point]

        # тут мы выгружаем выбранные точки по айдишникам
        #   топ-N точек с их embedding и временами
        try:
            cursor = connection.cursor()

            cursor.execute(f'''
            SELECT id, name, longitude, latitude, address, description FROM locations
            WHERE id IN ({" ,".join([str(x) for x in locations['id']])});
            ''')
            test_data = cursor.fetchall()

            cursor.close()
        except:
            self.logger.warning("LOCATIONS' INFO EXTRACTIONS PROBLEM")
    
        self.relevance_scores = {}
        self.relevance_scores[0] = 0

        for i, (loc_id, name, lon, lat, addr, description) in enumerate(test_data):
            newLocation = Location(
                loc_id, name, float(lon), float(lat), addr, description, None, None)
            self.all_nodes.append(newLocation)
            # print(user_point.get_latitude(), user_point.get_longitude(),
            #       newLocation.get_latitude(), newLocation.get_longitude())
            locations["time_to_user"][i] = \
                self.get_walking_time_from_and_to_user_in_minute(
                user_point, newLocation)
            self.relevance_scores[i+1] = locations['mark'][i] \
                * weight_rerank + (1 - weight_rerank) \
                * np.exp(-weight_distance *
                         (30 + locations["time_to_user"][i]) / T_limit)

        # print([x.get_title() for x in all_nodes])

        self.num_locations = len(self.all_nodes)

        # подгружаем время в минутах

        # 2. Создаем урезанную матрицу времени только для этих точек
        # Это нужно для эффективности работы callback-функции.

        self.logger.debug("СОЗДАНИЕ МАТРИЦЫ ВРЕМЕНИ ДЛЯ ВЫБРАННЫХ ТОЧЕК")

        self.small_time_matrix = []
        for from_node in self.all_nodes:
            row = []
            for to_node in self.all_nodes:
                # если точка - пользователь, то пользуемся OSRM
                if (from_node == user_point and to_node == user_point):
                    row.append(1)
                elif (from_node == user_point):
                    row.append(
                        int(locations["time_to_user"][locations["id"] == to_node.get_id()][0]) + 31)
                elif (to_node == user_point):
                    row.append(
                        int(locations["time_to_user"][locations["id"] == from_node.get_id()][0]) + 31)
                else:
                    row.append(self.loadTime(from_node.id, to_node.id))
            self.small_time_matrix.append(row)

        # Создаем менеджер индексов.
        # Стартуем с индекса 0 (user_start_index в small_time_matrix).
        self.manager = pywrapcp.RoutingIndexManager(
            self.num_locations,  # количество точек
            1,              # количество пешеходов
            0               # индекс стартовой точки в small_time_matrix
        )

        # Создаем модель маршрутизации
        self.routing = pywrapcp.RoutingModel(self.manager)

        # Определяем и регистрируем callback для времени
        def time_callback(from_index, to_index):
            from_node = self.manager.IndexToNode(from_index)
            to_node = self.manager.IndexToNode(to_index)

            return self.small_time_matrix[from_node][to_node]

        self.transit_callback_index = self.routing.RegisterTransitCallback(
            time_callback)

        # 4. Устанавливаем стоимость дуг равной времени
        self.routing.SetArcCostEvaluatorOfAllVehicles(
            self.transit_callback_index)

        # 5. Добавляем ограничение по времени
        time_dimension_name = 'Time'
        self.routing.AddDimension(
            self.transit_callback_index,
            0,               # slack_max:
            # capacity: общее ограничение
            #  по времени (например, 7200 секунд = 2 часа)
            int(self.time_limit),
            True,            # fix_start_cumul_to_zero: начинаем с времени 0
            time_dimension_name
        )
        self.time_dimension = self.routing.GetDimensionOrDie(
            time_dimension_name)

        MAX_PENALTY = 10**6
        # вознаграждения у нас - это наш релеванс,
        #  с увеличенным кэфом, чтобы больше ревенса собирало
        for node in range(1, self.num_locations):
            reward = int(self.relevance_scores[node] * 10000)
            self.routing.AddDisjunction(
                [self.manager.NodeToIndex(node)],
                MAX_PENALTY - reward  # чем больше reward, тем меньше penalty
            )
        self.routing.CloseModel()
        # настройки решателя
        # 8. Настраиваем стратегию поиска для задачи с вознаграждениями
        # настройки решателя (поиск маршрута с наградами)
        self.search_parameters = pywrapcp.DefaultRoutingSearchParameters()

        # стратегия начального решения
        self.search_parameters.first_solution_strategy = (
            routing_enums_pb2.FirstSolutionStrategy.PATH_CHEAPEST_ARC
        )

        # улучшение решения с учётом локальных перемещений
        self.search_parameters.local_search_metaheuristic = (
            routing_enums_pb2.LocalSearchMetaheuristic.GUIDED_LOCAL_SEARCH
        )

        # лимит по времени поиска
        self.search_parameters.time_limit.FromSeconds(5)

        # ограничиваем количество найденных решений (чтобы не застрять)
        self.search_parameters.solution_limit = 1

        # search_parameters.log_search = True -- ЛОГИ

    def get_walking_time_from_and_to_user_in_minute(
            self, FROMpoint: Location, TOpoint: Location):
        """Получает время пешего маршрута от OSRM"""
        try:
            from_lon, from_lat = FROMpoint.get_longitude(), FROMpoint.get_latitude()
            to_lon, to_lat = TOpoint.get_longitude(), TOpoint.get_latitude()

            url = f"http://localhost:8000/route/v1/foot/{from_lon},{from_lat};{to_lon},{to_lat}"
            response = requests.get(url, timeout=30)
            data = response.json()

            if data['code'] == 'Ok':
                duration_seconds = data['routes'][0]['duration']
                return int(duration_seconds / 60) + 1
        except Exception as e:
            self.logger.warning(f"OSRM PROBLEM --- {e}")
        return -1

    def solve(self):

        response = {
            "time": int,
            # время, которое пользователь затратит на всеь маршрут
            "description": "",
            # описание маршрута
            "count_places": int,
            # кол-во точек на машруте
            "places": []
        }

        # Решаем задачу!
        try:
            solution = self.routing.SolveWithParameters(self.search_parameters)
        except:
            self.logger.warning("SOLVE PROBLEM")

        if solution:
            # начинаем с первого (и единственного) транспортного средства
            index = self.routing.Start(0)
            plan_output = 'Маршрут:\n'
            route_nodes = []
            # Список для хранения индексов точек в порядке маршрута
            total_time = 0
            total_relevance = 0
            visited_nodes = set()

            while not self.routing.IsEnd(index):
                # Получаем индекс текущей точки в нашей small_time_matrix
                node_index = self.manager.IndexToNode(index)
                route_nodes.append(self.all_nodes[node_index])
                visited_nodes.add(node_index)
                response["places"].append({
                    "title": self.all_nodes[node_index].get_title(),
                    "addres": self.all_nodes[node_index].get_addr(),
                    "coordinate": [self.all_nodes[node_index].get_latitude(), self.all_nodes[node_index].get_longitude()],
                    "url": "",
                    "time_to_visit": 30 if node_index != 0 else 0,
                    "time_to_come": self.small_time_matrix[node_index][previous_index] - 30 if node_index != 0 else 0 ,
                    # время, чтобы добраться до места
                    "description": self.all_nodes[node_index].get_description()
                    # информация о месте
                })
                
                # Считаем relevance для посещенной точки (кроме стартовой) 
                if node_index != 0:
                    total_relevance += self.relevance_scores.get(node_index, 0)

                total_time += self.small_time_matrix[node_index][previous_index] - 30 if node_index != 0 else 0

                plan_output += f'   ({self.all_nodes[node_index].get_title()} \
                {self.all_nodes[node_index].get_id()}) ->'
                previous_index = node_index
                index = solution.Value(self.routing.NextVar(index))
                # Считаем общее время
                
            # Добавляем последнюю точку
            node_index = self.manager.IndexToNode(index)
            if node_index not in visited_nodes:  # Чтобы не дублировать подсчет
                route_nodes.append(self.all_nodes[node_index])
                if node_index != 0:
                    total_relevance += self.relevance_scores.get(node_index, 0)

            self.logger.debug(f"SOLVE : {plan_output}")

            plan_output += ' КОНЕЦ\n'

            plan_output += f'Общее время пути: {total_time - self.get_walking_time_from_and_to_user_in_minute(
                route_nodes[-1], route_nodes[0])} минут или же \
            {(total_time - self.get_walking_time_from_and_to_user_in_minute(
                    route_nodes[-1], route_nodes[0]))/60} часов.\n'
            plan_output += f'Суммарный relevance: {total_relevance}\n'
            plan_output += f'Посещено точек: {len(visited_nodes)} из \
            {self.num_locations}\n'

            response["time"] = (total_time)  
            response["description"] = ""
            response["count_places"] = len(visited_nodes)

            # print(plan_output)

        else:
            print('Решение не найдено.')
            self.logger.debug("SOLVE : РЕШЕНИЕ НЕ НАЙДЕНО")

        return response


# Создаем экземпляр класса и запускаем решение