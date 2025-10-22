
use("task_service_db");
db.createCollection("tasks")

db.tasks.createIndex({ userdId: 1, dueDate: 1});
db.tasks.createIndex({ userId: 1, status: 1, dueDate: 1});
db.tasks.createIndex({ "title": "text", "description": "text"});

db.tasks.getIndexes();