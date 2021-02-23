FROM python:3.8 as builder

COPY requirements.txt /requirements.txt
RUN mkdir /pythonlibs && pip install --target /pythonlibs -r /requirements.txt

FROM public.ecr.aws/lambda/python:3.8

COPY app.py ${LAMBDA_TASK_ROOT}
COPY --from=builder /pythonlibs ${LAMBDA_RUNTIME_DIR}
CMD ["app.handler"] 