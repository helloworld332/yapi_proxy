#!/usr/bin/env python
# -*- coding:utf-8 -*-

import sys
reload(sys)
sys.setdefaultencoding("utf-8")
from pymongo import *

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
  
  db.project.delete_one({"_id":project_id})
  db.interface_cat.delete_many({"project_id":project_id})
  db.interface.delete_many({"project_id":project_id})
  db.interface_case.delete_many({"project_id":project_id})
  db.interface_col.delete_many({"project_id":project_id})
