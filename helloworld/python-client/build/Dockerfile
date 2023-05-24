FROM python:3.11-slim
RUN mkdir /app
WORKDIR /app
ADD . /app
RUN pip config set global.index-url https://pypi.tuna.tsinghua.edu.cn/simple
RUN pip config set global.extra-index-url "https://mirrors.aliyun.com/pypi/simple https://mirrors.bfsu.edu.cn/pypi/web/simple"
RUN pip install --upgrade pip
RUN pip install -r requirements.txt
EXPOSE 5000
ENTRYPOINT ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "5000"]