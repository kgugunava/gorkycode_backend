import numpy as np
import threading

class Retriever:

    def __init__(self, connection, model, reranker):
        self.connection = connection
        self.model = model
        self.reranker = reranker
        self.encode_lock = threading.Lock()

    def get_embeddings_from_query(self, query):
        with self.encode_lock:
            query_embedding = self.model.encode(
                [query], batch_size=16, return_dense=True, return_sparse=True,
                return_colbert_vecs=True, convert_to_numpy=True)
        return query_embedding

    def get_reranker(self, query, top_50):

        pairs = [[query, top_50['description'][i]] for i in range(50)]
        reranker_score = self.reranker.compute_score(pairs,batch_size=16  )
        for i in range(50):
            top_50[i]['mark'] = 1 / (1 + np.exp(-reranker_score[i]))
        top_30 = np.sort(top_50, order='mark')[::-1][:30]
        return top_30

    def get_retriever(self, query_description, k_dense, k_sparse, k_colbert):
        cursor = self.connection.cursor()
        cursor.execute(
            "SELECT id, name, description, dense, sparse,"
            " colbert FROM locations")
        record = cursor.fetchall()
        cursor.close()
        dtype = np.dtype([('id', np.int32), ('name', 'U100'),
                         ('description', object), ('mark', np.float32),
                         ('time_to_user', np.int32)])
        ans = np.array([], dtype=dtype)
        query_embeddings = self.get_embeddings_from_query(query_description)
        for id, name, description, dense, sparse, colbert in record:
            dense = dense.strip("[]")
            dense = np.fromstring(dense, sep=',')
            colbert = np.array(colbert, dtype=np.float32)
            mark = np.array([(id, name, description, (
                query_embeddings["dense_vecs"][0] @ dense) * k_dense
                + self.model.compute_lexical_matching_score(
                query_embeddings["lexical_weights"][0], sparse) * k_sparse
                + self.model.colbert_score(query_embeddings["colbert_vecs"][0],
                                           colbert) * k_colbert, 0)],
                dtype=dtype)
            ans = np.append(ans, mark)
        top_50 = np.sort(ans, order='mark')[::-1][:50]
        return top_50

    def get_top_by_embeddings(self, query_description,
                              k_dense, k_sparse, k_colbert):
        sort_by_retriever = self.get_retriever(
            query_description, k_dense, k_sparse, k_colbert)
        sort_by_rerank = self.get_reranker(
            query_description, sort_by_retriever)
        return sort_by_rerank[["id", "name", "mark", "time_to_user"]]
