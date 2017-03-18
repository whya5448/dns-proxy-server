angular.module("myApp", ["ngTable"]);

(function() {
		"use strict";
		
		angular.module('myApp')
		.filter('ipArray', function($filter) {
			return (ipArray, errorMessage) => {
				return ipArray.join('.');
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
				self.tableParams = new NgTableParams({}, {
						filterDelay: 0,
						getData: function(params) {
							return $http.get('/hostname').then(function(data) {
								params.total(data.data.hostnames.length); // recal. page nav controls
								console.debug('m=getData, length=%d', data.data.hostnames.length, data.data.hostnames);
								for(var i=0; i < data.data.hostnames.length; i++){
									data.data.hostnames[i].id = Math.ceil(Math.random() * new Date().getTime());
								}
								originalData = data.data.hostnames;
								return data.data.hostnames;
							}, function(err){
								console.error('m=getData, status=error', err);
							});
						}
				});

				self.cancel = cancel;
				self.del = del;
				self.save = save;
				self.saveNewLine = saveNewLine;

				function cancel(row, rowForm) {
						var originalRow = resetRow(row, rowForm);
						angular.extend(row, originalRow);
				}

				function del(row) {
					console.debug('m=del, hostname=%s', row.hostname)
					_.remove(originalData, function(item) {
							return row === item;
					});
					$http.delete('/hostname', row).then(function(data) {
						console.debug('m=del, status=scucess')
						self.tableParams.reload().then(function(data) {
							if (data.length === 0 && self.tableParams.total() > 0) {
									self.tableParams.page(self.tableParams.page() - 1);
									self.tableParams.reload();
							}
						});
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
					var originalRow = resetRow(row, rowForm);
					angular.extend(originalRow, row);

					$http.put('/hostname', row).then(function(data) {
						console.debug('m=save, status=scucess')
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
						data: line,
						headers: {
							'Content-Type': 'application/json'
						}
					}).then(function(data){
						console.debug('m=saveNewLine, status=success');
						$scope.line = null;
						self.tableParams.reload().then(function(data) {
								if (data.length === 0 && self.tableParams.total() > 0) {
										self.tableParams.page(self.tableParams.page() - 1);
										self.tableParams.reload();
								}
						});
					}, function(err){
						console.error('m=saveNewLine, status=error', err);
					});
				};
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
