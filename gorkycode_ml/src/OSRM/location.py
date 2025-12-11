class Location:
    def __init__(
        self, id: int, title: str, longitude: float, latitude: float, addr: str, description: str,
            sparse_emb: dict = None, dense_emb: list = None):        
        self.id = id
        self.title = title
        self.longitude = longitude
        self.latitude = latitude
        self.addr = addr
        self.description = description
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

    def get_addr(self) -> str:
        return self.addr

    def get_description(self) -> str:
        return self.description

    def get_sparse(self) -> dict:
        return self.sparse_emb

    def get_dense(self) -> list:
        return self.dense_emb
