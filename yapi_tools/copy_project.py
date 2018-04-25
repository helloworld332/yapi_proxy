#!/usr/bin/env python
# -*- coding:utf-8 -*-

import sys
reload(sys)
sys.setdefaultencoding("utf-8")
from pymongo import *

def new_id(table):
  max_id = 0
  for id_info in table.find({},{"_id":1}): 
    new_id = id_info["_id"]
    if new_id > max_id:
      max_id = new_id
  return max_id + 100

def new_project(db, project_id):
  project_table = db.project
  project = project_table.find_one({"_id":project_id})
  new_project_id = new_id(project_table)
  project["_id"] = new_project_id
  project["name"] = project["name"] + "_copy"
  project_table.insert_one(project)
  return new_project_id

def mapping_catid(db, old_project_id, new_project_id):
  catid_table = db.interface_cat
  catids = [cat for cat in catid_table.find({"project_id":old_project_id})]
  catid_map = {}
  for cat in catids:
    old_catid = cat["_id"]
    new_catid = new_id(catid_table)
    cat["_id"] = new_catid
    cat["project_id"] = new_project_id
    catid_map[old_catid] = new_catid
    catid_table.insert_one(cat)
  return catid_map

def mapping_colid(db, old_project_id, new_project_id):
  colid_table = db.interface_col
  colids = [col for col in colid_table.find({"project_id":old_project_id})]
  colid_map = {}
  for col in colids:
    old_colid = col["_id"]
    new_colid = new_id(colid_table)
    col["_id"] = new_colid
    col["project_id"] = new_project_id
    col["test_report"] = "{}"
    colid_map[old_colid] = new_colid
    colid_table.insert_one(col)
  return colid_map

def copy_interface(db, old_project_id, new_project_id, catid_map, colid_map):
  interface_table = db.interface
  interfaces = [interface for interface in interface_table.find({"project_id":old_project_id})]
  interfaceid_map = {}
  for interface in interfaces:
    new_interface_id = new_id(interface_table)
    interfaceid_map[interface["_id"]] = new_interface_id
    interface["_id"] = new_interface_id
    interface["project_id"] = new_project_id
    interface["catid"] = catid_map[interface["catid"]]
    interface_table.insert_one(interface)

  interface_case_table = db.interface_case
  interface_cases = [interface_case for interface_case in interface_case_table.find({"project_id":old_project_id})]
  for interface_case in interface_cases:
    new_interface_case_id = new_id(interface_case_table)
    interface_case["_id"] = new_interface_case_id
    interface_case["project_id"] = new_project_id
    interface_case["interface_id"] = interfaceid_map[interface_case["interface_id"]]
    interface_case["col_id"] = colid_map[interface_case["col_id"]]
    interface_case_table.insert_one(interface_case)

if __name__=='__main__':
  # argv[1]: project_id
  if len(sys.argv) < 2:
    print >>sys.stderr, "usage: %s project_id" % (sys.argv[0])
    sys.exit(1)
  project_id = int(sys.argv[1])
  client = MongoClient()
  # db = client.yapi2
  db = client.yapi
  num = db.project.find({"_id":project_id}).count()
  if num != 1:
    print >> sys.stderr, "project_id %d not exist." % (project_id)
    sys.exit(2)
  print "copy project_id: %d, project_name: %s" % (project_id, db.project.find({"_id":project_id})[0]["name"])

  new_project_id = new_project(db, project_id)
  print "create new project_id: %d" %(new_project_id)
  catid_map = mapping_catid(db, project_id, new_project_id)
  colid_map = mapping_colid(db, project_id, new_project_id)
  print "add catids: %s" %(catid_map)
  print "add colids: %s" %(colid_map)

  copy_interface(db, project_id, new_project_id, catid_map, colid_map)
