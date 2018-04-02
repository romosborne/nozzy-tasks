function init() {
}

var savedAuthKey = "";
var viewModel;
var gu;

function LoadProjects() {
    var projects = $.Deferred();

    callApi(
        {
            url: '/api/tasks',
            headers: { 'Authorization': 'Bearer: ' + savedAuthKey },
            dataType: 'JSON',
            success: function (data) {
                projects.resolve(data);
            }
        });

    return projects.promise();
}

function onSignIn(googleUser) {
    gu = googleUser;
    refreshAccessToken().done(function(){
        LoadProjects().done(function (projects) {
            console.log(projects);
            proj = $.map(projects, function (p) {
                tasks = $.map(p.tasks, function (t) {
                    return new Task(t.id, t.title, t.completed)
                })
                return new Project(p.id, p.name, tasks)
            });

            $.each(proj, function (i, v) {
                viewModel.addExistingProject(v);
            });

            viewModel.signedIn(true);
        });
    });
}

function refreshAccessToken() {
    var prom = $.Deferred();
    var id_token = gu.getAuthResponse().id_token;
    console.log("ID Token: " + id_token);

    $.ajax({
        url: '/get_token',
        headers: { 'Authorization': 'Bearer: ' + id_token },
        success: function (authKey) {
            savedAuthKey = authKey;
            prom.resolve();
        }
    });

    return prom.promise()
}

function callApi(request){
    console.log("Beep");
    request.headers = {'Authorization': 'Bearer: ' + savedAuthKey };
    if (!request.statusCode){
        request.statusCode = {};
    }

    request.statusCode[401] = function(){
        refreshAccessToken().done(function() {
            callApi(request);
        });
    }
    $.ajax(request);
}

$(function () {
    viewModel = new ProjectsViewModel();
    ko.applyBindings(viewModel);
})

var Task = function (id, title, completed) {
    var self = this;
    self.id = id;
    self.title = title;
    self.completed = ko.observable(completed)

    self.toggleComplete = function () {
        callApi({
            url: '/api/tasks/completion',
            method: 'POST',
            data: JSON.stringify({
                "task_id": self.id,
                "completed": self.completed()
            })
        });

        return true;
    }
}

var Project = function (id, name, tasks) {
    var self = this;
    self.id = id;
    self.name = name;
    self.newTaskName = ko.observable("");
    self.tasks = ko.observableArray(tasks);

    self.enterAdd = function (d, e) {
        e.keyCode === 13 && self.addTaskInline();
        return true;
    };

    self.addTask = function (task) {
        self.tasks.push(task);
    };

    self.addTaskInline = function () {
        callApi({
            url: '/api/tasks',
            method: 'POST',
            data: JSON.stringify({
                "title": self.newTaskName(),
                "project_id": self.id
            }),
            success: function (response) {
                self.tasks.push(new Task(response.id, response.title, false));
                self.newTaskName("");
            }
        });
    }

    self.deleteTask = function (task) {
        callApi({
            url: '/api/tasks/' + task.id,
            method: 'DELETE',
        })

        self.tasks.remove(task);

        return true
    }
}

function ProjectsViewModel(projects) {
    var self = this;

    self.signedIn = ko.observable(false)

    self.projects = ko.observableArray(projects);
    self.newProjectName = ko.observable("")

    self.enterAdd = function (d, e) {
        e.keyCode === 13 && self.addProject();
        return true;
    };

    self.addProject = function () {
        callApi({
            url: '/api/project',
            method: 'POST',
            data: JSON.stringify({ "name": self.newProjectName() }),
            success: function (response) {
                console.log(response);
                self.projects.push(new Project(response.id, response.name, []));
                self.newProjectName("");
            }
        });
    }

    self.addExistingProject = function (project) {
        self.projects.push(project);
    }
}
