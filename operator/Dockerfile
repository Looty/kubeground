FROM python:3.12
RUN pip install kopf
ADD src/ /src
RUN pip install -r /src/requirements.txt
CMD kopf run /src/op.py --verbose
