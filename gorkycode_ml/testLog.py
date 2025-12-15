import logging
import sys
from typing import List, Optional

# Настройка логирования
def setup_logging(debug_mode: bool = False):
    """Настройка системы логирования"""
    # Создаем логгер
    logger = logging.getLogger()
    logger.setLevel(logging.DEBUG if debug_mode else logging.INFO)
    
    # Создаем форматтер
    formatter = logging.Formatter(
        '%(asctime)s - %(name)s - %(levelname)s - %(message)s',
        datefmt='%Y-%m-%d %H:%M:%S'
    )
    
    # Обработчик для консоли (stdout)
    console_handler = logging.StreamHandler(sys.stdout)
    console_handler.setLevel(logging.DEBUG if debug_mode else logging.INFO)
    console_handler.setFormatter(formatter)
    
    # Обработчик для файла (только INFO и выше)
    file_handler = logging.FileHandler('app.log', encoding='utf-8')
    file_handler.setLevel(logging.INFO)
    file_handler.setFormatter(formatter)
    
    # Добавляем обработчики к логгеру
    logger.addHandler(console_handler)
    logger.addHandler(file_handler)
    
    return logger

# Пример класса с логированием
class Calculator:
    """Класс для математических операций с логированием"""
    
    def __init__(self, name: str):
        self.logger = logging.getLogger(f'Calculator.{name}')
        self.logger.info(f'Инициализация калькулятора "{name}"')
        self.name = name
    
    def add(self, numbers: List[float]) -> float:
        """Сложение чисел"""
        self.logger.debug(f'Начало сложения: {numbers}')
        
        if not numbers:
            self.logger.warning('Пустой список чисел для сложения')
            return 0.0
        
        result = sum(numbers)
        self.logger.info(f'Сложение завершено: {numbers} = {result}')
        self.logger.debug(f'Детали сложения: {len(numbers)} чисел, сумма = {result}')
        
        return result
    
    def divide(self, a: float, b: float) -> Optional[float]:
        """Деление двух чисел"""
        self.logger.debug(f'Начало деления: {a} / {b}')
        
        if b == 0:
            self.logger.error(f'Попытка деления на ноль: {a} / {b}')
            return None
        
        result = a / b
        self.logger.info(f'Деление завершено: {a} / {b} = {result}')
        self.logger.debug(f'Детали деления: тип результата = {type(result)}')
        
        return result
    
    def calculate_average(self, numbers: List[float]) -> Optional[float]:
        """Вычисление среднего значения"""
        self.logger.debug(f'Вычисление среднего для списка: {numbers}')
        
        if not numbers:
            self.logger.warning('Пустой список для вычисления среднего')
            return None
        
        try:
            total = self.add(numbers)
            count = len(numbers)
            average = total / count
            
            self.logger.info(f'Среднее значение: {numbers} = {average}')
            self.logger.debug(f'Детали: сумма = {total}, количество = {count}')
            
            return average
        except Exception as e:
            self.logger.error(f'Ошибка при вычислении среднего: {e}', exc_info=True)
            return None

# Основная функция программы
def main():
    """Основная функция программы"""
    # Настройка логирования с DEBUG уровнем
    logger = setup_logging(debug_mode=True)
    logger.info('=' * 50)
    logger.info('Запуск программы')
    logger.info('=' * 50)
    
    # Создаем экземпляр калькулятора
    logger.debug('Создание экземпляра калькулятора')
    calc = Calculator("Основной")
    testCalc = Calculator("ТЕСТОВЫЙ")
    
    # Примеры операций с разным уровнем логирования
    logger.info('Начало выполнения математических операций')
    
    # 1. Сложение
    numbers_to_add = [1.5, 2.5, 3.5, 4.5]
    logger.debug(f'Подготовка данных для сложения: {numbers_to_add}')
    sum_result = calc.add(numbers_to_add)
    logger.info(f'Результат сложения: {sum_result}')
    
    # 2. Деление
    logger.debug('Подготовка данных для деления')
    division_result = testCalc.divide(10, 2)
    if division_result is not None:
        logger.info(f'Результат деления: {division_result}')
    
    # 3. Деление на ноль (ошибка)
    logger.debug('Попытка деления на ноль')
    error_division = calc.divide(10, 0)
    if error_division is None:
        logger.warning('Деление на ноль предотвращено')
    
    # 4. Вычисление среднего
    logger.debug('Вычисление среднего значения')
    average_numbers = [10, 20, 30, 40, 50]
    average_result = calc.calculate_average(average_numbers)
    if average_result is not None:
        logger.info(f'Среднее значение: {average_result}')
    
    # 5. Ошибка с пустым списком
    logger.debug('Попытка вычисления среднего для пустого списка')
    empty_average = calc.calculate_average([])
    if empty_average is None:
        logger.info('Среднее для пустого списка не вычислено')
    
    # Завершение программы
    logger.info('=' * 50)
    logger.info('Программа завершена успешно')
    logger.info('=' * 50)
    logger.debug('Освобождение ресурсов...')

if __name__ == "__main__":
    main()