angular.module("myApp", ["ngTable"]);

(function() {
		"use strict";
		
		angular.module('myApp')
		.filter('ipArray', function($filter) {
			return (ipArray, errorMessage) => {
				return ipArray.join('.');
			}
		})
		.controller("demoController", ["NgTableParams", "$http", demoController])
		.directive('arrayToString', function() {
			return {
				require: 'ngModel',
				link: function(scope, element, attrs, ngModel) {
					ngModel.$parsers.push(function(value) {
						return value.split('\\.');
					});
					ngModel.$formatters.push(function(value) {
						return value.join('.');
					});
				}
			};
		});;

		function demoController(NgTableParams, $http) {
				var self = this;

				var originalData = dataFactory();

				self.tableParams = new NgTableParams({}, {
						filterDelay: 0,
//						dataset: angular.copy(originalData),
						getData: function(params) {
							console.debug('loading data');
							// ajax request to api
							return $http.get('/hostname').then(function(data) {
								params.total(data.inlineCount); // recal. page nav controls
								console.debug('loaded data', data);
								return data.data.hostnames;
							}, function(err){
								console.debug('err', err);
							});
						}
				});

				self.cancel = cancel;
				self.del = del;
				self.save = save;
				self.saveNewLine = saveNewLine;

				//////////


				function cancel(row, rowForm) {
						var originalRow = resetRow(row, rowForm);
						angular.extend(row, originalRow);
				}

				function del(row) {
						_.remove(self.tableParams.settings().dataset, function(item) {
								return row === item;
						});
						self.tableParams.reload().then(function(data) {
								if (data.length === 0 && self.tableParams.total() > 0) {
										self.tableParams.page(self.tableParams.page() - 1);
										self.tableParams.reload();
								}
						});
				}

				function resetRow(row, rowForm){
						row.isEditing = false;
						rowForm.$setPristine();
						//self.tableTracker.untrack(row);
						return _.findWhere(originalData, function(r){
								return r.id === row.id;
						});
				}

				function save(row, rowForm) {
						var originalRow = resetRow(row, rowForm);
						angular.extend(originalRow, row);
				}

				function saveNewLine(line){
						line = angular.copy(line);

						$http({method: 'POST', url: '/hostname/new/', data: line, headers: {'Content-Type': 'application/json'}}).then(function(data){
							console.debug('success', data);
						}, function(err){
							console.debug('err', err);
						});

						line.id = new Date().getTime();
						console.debug('m=saveNewLine, line=%o', line, self.tableParams, self.tableParams.settings());
//						self.tableParams.settings().dataset.push(line);
						self.tableParams.reload().then(function(data) {
								if (data.length === 0 && self.tableParams.total() > 0) {
										self.tableParams.page(self.tableParams.page() - 1);
										self.tableParams.reload();
								}
						});
				}
		}
})();



(function() {
		"use strict";

		angular.module("myApp").run(configureDefaults);
		configureDefaults.$inject = ["ngTableDefaults"];

		function configureDefaults(ngTableDefaults) {
				ngTableDefaults.params.count = 5;
				ngTableDefaults.settings.counts = [];
		}
})();



function  dataFactory() {
		return [{"id":1,"name":"Nissim","age":41,"money":454},{"id":2,"name":"Mariko","age":10,"money":-100},{"id":3,"name":"Mark","age":39,"money":291},{"id":4,"name":"Allen","age":85,"money":871},{"id":5,"name":"Dustin","age":10,"money":378},{"id":6,"name":"Macon","age":9,"money":128},{"id":7,"name":"Ezra","age":78,"money":11},{"id":8,"name":"Fiona","age":87,"money":285},{"id":9,"name":"Ira","age":7,"money":816},{"id":10,"name":"Barbara","age":46,"money":44},{"id":11,"name":"Lydia","age":56,"money":494},{"id":12,"name":"Carlos","age":80,"money":193}];
}