import configparser
import requests
import json
from flask import *
from datetime import datetime,timedelta

config = configparser.RawConfigParser()
config.read('/etc/statuzpage-frontend/config.cfg')

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

@app.route("/")
def index():
  incidentsData = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/incidents", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
  incidentsClosed = json.loads(requests.get("http://"+str(config.get('conf','apiHost'))+"/incidentsclosed", headers={"statuzpage-token":str(config.get('conf','apiToken'))}).content)
  if incidentsData == None:
    totalOpen = 0
  else:
    totalOpen = len(incidentsData)
  if incidentsClosed == None:
    totalClosed = 0
  else:
    totalClosed = len(incidentsClosed)    
  return render_template('index.html',lastincidents=incidentsData,incidentsclosed=incidentsClosed,totalOpen=totalOpen,totalClosed=totalClosed)

if __name__ == '__main__':
  app.run(host=str(config.get('conf','ip')),port=int(config.get('conf','port')),debug=True)