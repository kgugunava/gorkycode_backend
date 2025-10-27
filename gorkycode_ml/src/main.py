import psycopg2
from FlagEmbedding import BGEM3FlagModel, FlagReranker
from yaml import safe_load
# import json
# import pandas as pd

from RETRIEVER.retriever import Retriever
# from MODEL.model import Model
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
    with open('src/config/config.yml', 'r') as f:
        data = safe_load(f)
    return data


def load_model():
    model = BGEM3FlagModel(data["MODEL"]["name"], local_files_only=True,
                           use_fp16=True)
    return model


def load_reranker():
    reranker = FlagReranker(
        data["MODEL"]["reranker"], local_files_only=True, use_fp16=True)
    return reranker


data = load_config()
connection = connect_database(
    data["DB"]["user"], data["DB"]["password"], data["DB"]["host"],
    data["DB"]["port"], data["DB"]["database"])


query = {
    "interests": "Я хочу посмотреть на живописные виды. Посмотреть на объекты истории. Побывать в парках, а также в самых известных метсах",
    "time": 200,
    "coordinates": [43.853400, 56.261660,],
    "user_info": ""
}
load_model()

retriever = Retriever(connection, load_model(), load_reranker())
locations = retriever.get_top_by_embeddings(
    query["interests"], data["BGE"]["dense"],
    data["BGE"]["sparse"], data["BGE"]["colbert"])
print(locations[["name", "mark"]])
userPonit = Location(
    0, "UserPoint", query["coordinates"][0], query["coordinates"][1])
planner = RoutePlanner(connection, userPonit, query["time"],
                       locations, data["RELEVANCE"]["weight_reranker"],
                       data["RELEVANCE"]["weight_distance"])
planner.solve()

connection.close()
