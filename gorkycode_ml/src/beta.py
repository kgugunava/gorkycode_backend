from vllm import LLM, SamplingParams

# 0.5B-вариант влезет в 4 ГБ VRAM
model_name = "Qwen/Qwen2.5-0.5B-Instruct"

# Настройка генерации
params = SamplingParams(
    temperature=0.6,
    top_p=0.9,
    max_tokens=150
)

llm = LLM(
    model=model_name,
    tensor_parallel_size=1,  # одна GPU
    dtype="float16"
)

prompt = """Ты — эксперт по культурному и туристическому планированию маршрутов.
На вход получаешь JSON маршрута и запрос пользователя.
Опиши коротко (3–6 предложений), почему маршрут составлен именно так."""

outputs = llm.generate(prompt, sampling_params=params)

print("\n🧭 Ответ:\n", outputs[0].outputs[0].text)
