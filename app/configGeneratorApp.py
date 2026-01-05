import os
from fastapi import FastAPI
from pydantic import BaseModel
from jinja2 import Environment, FileSystemLoader
import re

app = FastAPI()
env = Environment(loader=FileSystemLoader("templates"))

class templateRequest(BaseModel):
    config_type: str

class renderRequest(BaseModel):
    config_type: str
    data: dict

@app.get("/v1/health")
def healthCheckHandler():
    return {"status":"ok"}

@app.post("/v1/render") 
def renderConfig(request: renderRequest):
    templateName = f"{request.config_type}.j2"
    template = env.get_template(templateName)
    config = template.render(**request.data)

    filename = f"{request.config_type}.txt"

    return {
        "config": config,
        "filename": filename
    }

@app.get("/v1/templates")
def listTemplates():
    templates = []
    for file in os.listdir("templates"):
        if file.endswith(".j2"):
            templates.append({"name":(file.replace(".j2", ""))})
    return {"templates": templates}

@app.post("/v1/getTemplate")
def getTemplate(request: templateRequest):
    template_file = f"{request.config_type}.j2"
    file_path = os.path.join("templates", template_file)

    if not os.path.exists(file_path):
        return {"fields": []}

    with open(file_path, "r", encoding="utf-8") as f:
        source = f.read()

    comment_matches = re.findall(
        r"\{\#-?\s*(\d+)\s*:\s*([\w]+)\s*:\s*([^:]+?)(?:\s*:\s*(.+?))?\s*-?\#\}",
        source
    )

    ordered_fields = []
    for order, field, desc, regex in sorted(comment_matches, key=lambda x: int(x[0])):
        pattern = regex.strip() if regex else ""

        pattern = pattern.replace(r"\{", "{").replace(r"\}", "}")

        ordered_fields.append({
            "name": field.strip(),
            "description": desc.strip(),
            "pattern": pattern
        })

    return {"fields": ordered_fields}