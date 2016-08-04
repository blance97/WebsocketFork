var app = angular.module("chat", []);

app.controller("MainCt1", ["$scope", function($scope) {
    $scope.messages = []; // list that we recive from server

    var conn = new WebSocket("ws://localhost/ws");

    conn.onclose = function(event) {
        $scope.$apply(function() {
            $scope.messages.push("DISCONNECTED")
        })
    }
    conn.onopen = function(event) {
        $scope.$apply(function() {
            $scope.messages.push("CONNECTED")
        })
    }

    conn.onmessage = function(event) {
        $scope.$apply(function() {
            $scope.messages.push(event.data);
        })
    }
    $scope.send = function(){
      conn.send($scope.msg);
      $scope.msg = "";
    }

}])
