import torch
from packaging import version
import threading

from flask import Response
import json

import psycopg2

from FlagEmbedding import BGEM3FlagModel, FlagReranker
from flask import Flask, request, jsonify
from yaml import safe_load



# import json
# import pandas as pd

from RETRIEVER.retriever import Retriever
from MODEL.model import Model
from OSRM.location import Location
from OSRM.route_planner import RoutePlanner


def connect_database(user: str, password: str, host: str, port: str,
                     database: str):
    connection = psycopg2.connect(user=user,
                                  password=password,
                                  host=host,
                                  port=port,
                                  database=database)
    return connection


# def make_embeddings(connection, model, dir):
#     data = pd.read_csv(dir, sep=";")
#     cursor = connection.cursor()
#     corpus = list(data["description"])
#     print(len(corpus))
#     embedding = model.encode(corpus, return_dense=True,
#                              return_sparse=True, return_colbert_vecs=True,
#                              convert_to_numpy=True)
#     for i in range(259):
#         colbert_vec = json.dumps(embedding["colbert_vecs"][i].tolist())
#         sparse = {int(k): float(v)
#                   for k, v in dict(embedding["lexical_weights"][i]).items()}
#         json_sparse = json.dumps(sparse)
#         dense = list(embedding["dense_vecs"][i])
#         vector_dense = "[" + ",".join(map(str, dense)) + "]"
#         query = ''' UPDATE locations SET dense = %s::vector,
#         sparse = %s, colbert = %s WHERE id = %s;'''
#         params = (vector_dense,
#                   json_sparse, colbert_vec, int(data["id"][i]))
#         cursor.execute(query, params)
#     connection.commit()
#     cursor.close()


# def add_description(connection, dir):
#     data = pd.read_csv(dir, sep=";")
#     cursor = connection.cursor()
#     corpus = list(data["description"])
#     for i in range(259):
#         query = ''' UPDATE locations SET description = %s WHERE id = %s;'''
#         params = (corpus[i], int(data["id"][i]))
#         cursor.execute(query, params)
#     connection.commit()
#     cursor.close()


def load_config():
    with open('src/config/config.yml','r', encoding='utf-8') as f:
        data = safe_load(f)
    return data


def load_model(device):
    model = BGEM3FlagModel(data["MODEL"]["name"],device=device, local_files_only=True,
                           use_fp16=True)
    model.model.eval()
    return model


def load_reranker(device):
    reranker = FlagReranker(
        data["MODEL"]["reranker"],device=device, local_files_only=True, use_fp16=True)
    return reranker


data = load_config()


# load_model()


# query = {
#    "interests": "Я хочу посмотреть на живописные виды.
#  Посмотреть на объекты истории.
# Побывать в парках, а также в самых известных метсах",
#    "time": 200,
#    "coordinates": [43.853400, 56.261660,],
#    "user_info": ""
# }

# всё, что ниже - работа сервера
device = "cuda" if torch.cuda.is_available() else "cpu"
app = Flask(__name__)
modelx = load_model(device)
rerankerx = load_reranker(device)
query_lock = threading.Lock()
# query_lock = threading.Lock()

@app.route('/route', methods=['POST'])
def handle_route_request():
    with query_lock:
        connection = connect_database(
            data["DB"]["user"], data["DB"]["password"], data["DB"]["host"],
            data["DB"]["port"], data["DB"]["database"])
        # полуяам json из запроса
        jsonchik = request.get_json()
        retriever = Retriever(connection, modelx, rerankerx)
        # блок кода для обработки запроса моделькой
        locations = retriever.get_top_by_embeddings(
            jsonchik.get('interests'), data["BGE"]["dense"],
            data["BGE"]["sparse"], data["BGE"]["colbert"]
        )


        userPonit = Location(
            0, "UserPoint", jsonchik.get('coordinates')[1],
            jsonchik.get('coordinates')[0], "", "")
        planner = RoutePlanner(connection, userPonit, jsonchik.get('time_for_route'),
                            locations, data["RELEVANCE"]["weight_reranker"],
                            data["RELEVANCE"]["weight_distance"])

        locations = planner.solve()
        model = Model(data["OLLAMA"]["name"])
        model.request_to_model(data["OLLAMA"]["user_prompt"], locations, jsonchik.get('interests'))
        connection.close()
        # print(locations)
    return jsonify(locations)
# Response(json.dumps(locations, ensure_ascii=False, default=str), content_type="application/json")
   


# planner.solve()


# тут начинается область работы сервера
if __name__ == '__main__':
    app.run(debug=True, host='0.0.0.0', port=5001)
