(function () {
	console.log("Hello angular!")
	var app = angular.module("myApp", ["ngResource"]);
	app.factory("Mollusks", function ($resource) {
		return $resource("http://localhost\\:31415/mollusk/:id", {id: "@id"});
	})
	app.controller("HelloCtrl", function ($scope, Mollusks) {
		console.log(Mollusks)
		$scope.mollusks = Mollusks.query()

		$scope.addMollusk = function () {
			var newMollusk = new Mollusks($scope.newMollusk)
			newMollusk.$save()
			$scope.mollusks = Mollusk.$query()
			return false
		}

		$scope.appState = "ready!"
		console.log("HelloCtrl  ready!")
	})
})();
