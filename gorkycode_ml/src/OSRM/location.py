class Location:
    def __init__(
        self, id: int, title: str, longitude: float, latitude: float,
            sparse_emb: dict = None, dense_emb: list = None):
        self.id = id
        self.title = title
        self.longitude = longitude
        self.latitude = latitude
        self.sparse_emb = sparse_emb if sparse_emb is not None else {}
        self.dense_emb = dense_emb if dense_emb is not None else []

    def get_id(self) -> int:
        return self.id

    def get_title(self) -> str:
        return self.title

    def get_longitude(self) -> float:
        return self.longitude

    def get_latitude(self) -> float:
        return self.latitude

    def get_sparse(self) -> dict:
        return self.sparse_emb

    def get_dense(self) -> list:
        return self.dense_emb
