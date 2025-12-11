import numpy as np
import threading
import logging


class Retriever:

    def __init__(self, connection, model, reranker):

        self.logger = logging.getLogger("RETRIEVER")
        
        self.connection = connection
        self.model = model
        self.reranker = reranker
        self.encode_lock = threading.Lock()

    def get_embeddings_from_query(self, query):
        self.logger.debug("GET_EMBEDDINGS_FROM_QUERY IS LAUNCHED")
        self.logger.debug(f"GET_EMBEDDINGS_FROM_QUERY: INPUT QUERY: {query}")        
        try:
            with self.encode_lock:
                query_embedding = self.model.encode(
                    [query], batch_size=16, return_dense=True, return_sparse=True,
                    return_colbert_vecs=True, convert_to_numpy=True)
        except:
            self.logger.warning("GET_ENBEDDINGS_FROM_QUERY IS BROKEN")
        return query_embedding

    def get_reranker(self, query, top_50):
        self.logger.debug(f"GET RERANKER IS LAUNCHED")
        try:
            pairs = [[query, top_50['description'][i]] for i in range(50)]
            reranker_score = self.reranker.compute_score(pairs,batch_size=16  )
            for i in range(50):
                top_50[i]['mark'] = 1 / (1 + np.exp(-reranker_score[i]))
            top_30 = np.sort(top_50, order='mark')[::-1][:30]
        except:
            self.logger.warning("GET_RERANKER IS BROKEN")
    
        return top_30

    def get_retriever(self, query_description, k_dense, k_sparse, k_colbert):

        self.logger.debug("GET_RETRIEVER IS LAUNCHED")

        try:
            cursor = self.connection.cursor()
            cursor.execute(
                "SELECT id, name, description, dense, sparse,"
                " colbert FROM locations")
            record = cursor.fetchall()
            cursor.close()
        except:
            self.logger.warning("GET_RETRIEVER CONNECTION WITH DATABASE PROBLEM")

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
