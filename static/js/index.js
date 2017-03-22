angular.module("myApp", ["ngTable"]);

(function() {
		"use strict";
		
		angular.module('myApp')
		.filter('ipArray', function($filter) {
			return (ipArray, errorMessage) => {
				return ipArray.join('.');
			}
		})
		.filter('envFormatter', function($filter) {
			return (env, errorMessage) => {
				return env ? env : 'Default';
			}
		})
		.controller("demoController", ["NgTableParams", "$http", "$scope", demoController])
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
			}
		})
		.directive('ip', function() {
			return {
				require: 'ngModel',
				link: function(scope, elm, attrs, ctrl) {
					ctrl.$validators.ip = function(modelValue, viewValue) {

						if( ctrl.$isEmpty(modelValue)){
							return true;
						}

						if (/^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$/.test(viewValue)) {
							console.debug('m=ip, status=valid');
							return true;
						}

						console.debug('m=ip, status=INvalid');
						return false;
					}
				}
			}
		});

		function demoController(NgTableParams, $http, $scope) {
				var self = this;
				var originalData;
				$scope.activeEnv = "";
				$scope.envs = [];
				$scope.envs = [];
				$http.get('/env').then(function(data) {
					console.debug('m=getEnvs, length=%d', data.data.length, data.data);
					$scope.envs = data.data;
				}, function(err){
					console.error('m=getData, status=error', err);
				});

				$http.get('/env/active').then(function(data) {
					console.debug('m=getActiveEnv', data.data);
					$scope.activeEnv = data.data.name;
					reloadTable();
				}, function(err){
					console.error('m=getActiveEnv, status=error', err);
				});

				self.tableParams = new NgTableParams({}, {
						filterDelay: 0,
						getData: function(params) {
							return $http.get('/hostname/?env=' + $scope.activeEnv).then(function(data) {
								params.total(data.data.hostnames.length); // recal. page nav controls
								console.debug('m=getData, length=%d', data.data.hostnames.length, data.data.hostnames);
								originalData = data.data.hostnames;
								return angular.copy(data.data.hostnames);
							}, function(err){
								console.error('m=getData, status=error', err);
							});
						}
				});

				self.cancel = cancel;
				self.del = del;
				self.save = save;
				self.saveNewLine = saveNewLine;
				self.requestEdit = requestEdit;

				function cancel(row, rowForm) {
					console.debug('m=cancel, status=begin, row=%o', row);
					var originalRow = resetRow(row, rowForm);
					console.debug('m=cancel, status=success, foundRow=%o', originalRow);
					angular.extend(row, originalRow);
				}

				function requestEdit(row){
					console.debug('m=requestEdit, status=begin, row=%o', row);
					row.isEditing = true;
				}

				function del(row) {
					console.debug('m=del, hostname=%s', row.hostname)
					$http({
						url: '/hostname',
						method: 'DELETE',
						data: {env: $scope.activeEnv, hostname: row.hostname}
					}).then(function(data) {
						console.debug('m=del, status=scucess')
						_.remove(originalData, function(item) {
								return row.id === item.id;
						});
						reloadTable();
					}, function(err){
						console.error('m=save, status=error', err);

					});

				}

				function resetRow(row, rowForm){
						row.isEditing = false;
						rowForm.$setPristine();
						return _.findWhere(originalData, function(r){
								return r.id === row.id;
						});
				}

				function save(row, rowForm) {
					console.debug('m=save, hostname=%s', row.hostname)
					$http.put('/hostname', {env: $scope.activeEnv, hostname: row.hostname, ip: row.ip, ttl: ip.ttl}).then(function(data) {
						console.debug('m=save, status=scucess')
						var originalRow = resetRow(row, rowForm);
						angular.extend(originalRow, row);
					}, function(err){
						console.error('m=save, status=error', err);
					});

				}

				function saveNewLine(line, form){
					console.debug('m=saveNewLine, statys=begin, valid=%s', form.$valid)
					if(!form.$valid){
						return ;
					}
					line = angular.copy(line);
					console.debug('m=saveNewLine, status=begin, hostname=%o', line);
					line.ip = line.ip.split('\.').map(n => { return parseInt(n) });

					$http({
						method: 'POST', url: '/hostname/',
						data: {env: $scope.activeEnv, hostname: line.hostname, ip: line.ip, ttl: ip.ttl},
						headers: {
							'Content-Type': 'application/json'
						}
					}).then(function(data){
						console.debug('m=saveNewLine, status=success');
						$scope.line = null;
						reloadTable();
					}, function(err){
						console.error('m=saveNewLine, status=error', err);
					});
				};

				$scope.changeEnv = function(activeEnv){
					console.debug('m=changeEnv, status=begin, activeEnv=%o', activeEnv)
					$http.put('/env/active', {name: activeEnv}).then(function(data) {
						console.debug('m=changeEnv, status=scucess')
						reloadTable();
					}, function(err){
						console.error('m=changeEnv, status=error', err);
					});
				}

				function reloadTable(){
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
				ngTableDefaults.params.count = 500000;
				ngTableDefaults.settings.counts = [];
		}
})();
