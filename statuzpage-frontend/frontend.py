import configparser
import requests
import json
from flask import *
from datetime import datetime,timedelta

config = configparser.RawConfigParser()
config.read('/home/pia/GoogleDrive/Linux/statuzpage/statuzpage-frontend/config.cfg')

app = Flask(__name__)
app.secret_key = 'BYG>.L*((*$jj2h>#'

def duration(start,finish):
    timeStart = datetime.strptime(start, '%Y-%m-%d %H:%M:%S')
    timeFinish = datetime.strptime(finish, '%Y-%m-%d %H:%M:%S')
    duration = timeFinish-timeStart
    return duration

def onlyTime(date):
    data = datetime.strptime(date, '%Y-%m-%d %H:%M:%S')
    return data.strftime("%H:%M:%S")

def onlyDate(date):
    data = datetime.strptime(date, '%Y-%m-%d %H:%M:%S')
    return data.strftime("%d/%m/%Y")

def formatDate(date):
    data = datetime.strptime(date, '%Y-%m-%d %H:%M:%S')
    return data.strftime("%d/%m/%Y - %H:%M:%S")    

app.jinja_env.filters['onlytime'] = onlyTime
app.jinja_env.filters['onlydate'] = onlyDate
app.jinja_env.filters['formatdate'] = formatDate
app.jinja_env.globals.update(duration=duration)

@app.route("/tv")
def tv():
    totalOpen = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/totalOpen", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
    groupData = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/groups", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)    
    return render_template('tv.html',result=groupData,totalOpen=totalOpen["message"])

@app.route("/reports")
def reports():
    groupData = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/groups", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
    return render_template('reports.html',result=groupData)

@app.route("/reportsgroup", methods=['POST'])
def reportsgroup():
    if request.method == 'POST':
        incidentsData = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/incidents/"+request.form['idgroup']+"/"+request.form['month']+"/"+request.form['year'], headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
        group = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/group/"+request.form['idgroup'], headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
        return render_template('reportsgroup.html',result=incidentsData,group=group,month=request.form['month'],year=request.form['year'])
        # return jsonify(incidentsData)

@app.route("/")
def index():
    totalOpen = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/totalOpen", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
    groupData = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/groups", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
    incidentsData = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/incidents/limit/10", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)    
    return render_template('index.html',result=groupData,lastincidents=incidentsData,totalOpen=totalOpen["message"])

if __name__ == '__main__':
    app.run(host=str(config.get('conf','ip')),port=int(config.get('conf','port')))