angular.module("myApp", ["ngTable"]);

(function() {
	"use strict";

	angular.module("myApp").controller("demoController", demoController);
	demoController.$inject = ["NgTableParams"];


	function demoController(NgTableParams) {
		var self = this;

		var originalData = dataFactory();

		self.tableParams = new NgTableParams({}, {
			filterDelay: 0,
			dataset: angular.copy(originalData)
		});

		self.cancel = cancel;
		self.del = del;
		self.save = save;

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
			self.tableTracker.untrack(row);
			return _.findWhere(originalData, function(r){
				return r.id === row.id;
			});
		}

		function save(row, rowForm) {
			var originalRow = resetRow(row, rowForm);
			angular.extend(originalRow, row);
		}
	}
})();

(function() {
	"use strict";

	angular.module("myApp").controller("dynamicDemoController", dynamicDemoController);
	dynamicDemoController.$inject = ["NgTableParams"];

	function dynamicDemoController(NgTableParams) {
		var simpleList = dataFactory();
		var self = this;

		self.cols = [{
			field: "name",
			title: "Name",
			filter: {
				name: "text"
			},
			sortable: "name",
			dataType: "text"
		}, {
			field: "age",
			title: "Age",
			filter: {
				age: "number"
			},
			sortable: "age",
			dataType: "number"
		}, {
			field: "money",
			title: "Money",
			filter: {
				money: "number"
			},
			sortable: "money",
			dataType: "number"
		}, {
			field: "action",
			title: "",
			dataType: "command"
		}];

		var originalData = angular.copy(simpleList);

		self.tableParams = new NgTableParams({}, {
			filterDelay: 0,
			dataset: angular.copy(simpleList)
		});

		self.cancel = cancel;
		self.del = del;
		self.save = save;

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
			self.tableTracker.untrack(row);
			return _.findWhere(originalData, function(r){
				return r.id === row.id;
			});
		}

		function save(row, rowForm) {
			var originalRow = resetRow(row, rowForm);
			angular.extend(originalRow, row);
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

/**********
	The following directives are necessary in order to track dirty state and validity of the rows
	in the table as the user pages within the grid
------------------------
*/

(function() {
	angular.module("myApp").directive("demoTrackedTable", demoTrackedTable);

	demoTrackedTable.$inject = [];

	function demoTrackedTable() {
		return {
			restrict: "A",
			priority: -1,
			require: "ngForm",
			controller: demoTrackedTableController
		};
	}

	demoTrackedTableController.$inject = ["$scope", "$parse", "$attrs", "$element"];

	function demoTrackedTableController($scope, $parse, $attrs, $element) {
		var self = this;
		var tableForm = $element.controller("form");
		var dirtyCellsByRow = [];
		var invalidCellsByRow = [];

		init();

		////////

		function init() {
			var setter = $parse($attrs.demoTrackedTable).assign;
			setter($scope, self);
			$scope.$on("$destroy", function() {
				setter(null);
			});

			self.reset = reset;
			self.isCellDirty = isCellDirty;
			self.setCellDirty = setCellDirty;
			self.setCellInvalid = setCellInvalid;
			self.untrack = untrack;
		}

		function getCellsForRow(row, cellsByRow) {
			return _.find(cellsByRow, function(entry) {
				return entry.row === row;
			})
		}

		function isCellDirty(row, cell) {
			var rowCells = getCellsForRow(row, dirtyCellsByRow);
			return rowCells && rowCells.cells.indexOf(cell) !== -1;
		}

		function reset() {
			dirtyCellsByRow = [];
			invalidCellsByRow = [];
			setInvalid(false);
		}

		function setCellDirty(row, cell, isDirty) {
			setCellStatus(row, cell, isDirty, dirtyCellsByRow);
		}

		function setCellInvalid(row, cell, isInvalid) {
			setCellStatus(row, cell, isInvalid, invalidCellsByRow);
			setInvalid(invalidCellsByRow.length > 0);
		}

		function setCellStatus(row, cell, value, cellsByRow) {
			var rowCells = getCellsForRow(row, cellsByRow);
			if (!rowCells && !value) {
				return;
			}

			if (value) {
				if (!rowCells) {
					rowCells = {
						row: row,
						cells: []
					};
					cellsByRow.push(rowCells);
				}
				if (rowCells.cells.indexOf(cell) === -1) {
					rowCells.cells.push(cell);
				}
			} else {
				_.remove(rowCells.cells, function(item) {
					return cell === item;
				});
				if (rowCells.cells.length === 0) {
					_.remove(cellsByRow, function(item) {
						return rowCells === item;
					});
				}
			}
		}

		function setInvalid(isInvalid) {
			self.$invalid = isInvalid;
			self.$valid = !isInvalid;
		}

		function untrack(row) {
			_.remove(invalidCellsByRow, function(item) {
				return item.row === row;
			});
			_.remove(dirtyCellsByRow, function(item) {
				return item.row === row;
			});
			setInvalid(invalidCellsByRow.length > 0);
		}
	}
})();

(function() {
	angular.module("myApp").directive("demoTrackedTableRow", demoTrackedTableRow);

	demoTrackedTableRow.$inject = [];

	function demoTrackedTableRow() {
		return {
			restrict: "A",
			priority: -1,
			require: ["^demoTrackedTable", "ngForm"],
			controller: demoTrackedTableRowController
		};
	}

	demoTrackedTableRowController.$inject = ["$attrs", "$element", "$parse", "$scope"];

	function demoTrackedTableRowController($attrs, $element, $parse, $scope) {
		var self = this;
		var row = $parse($attrs.demoTrackedTableRow)($scope);
		var rowFormCtrl = $element.controller("form");
		var trackedTableCtrl = $element.controller("demoTrackedTable");

		self.isCellDirty = isCellDirty;
		self.setCellDirty = setCellDirty;
		self.setCellInvalid = setCellInvalid;

		function isCellDirty(cell) {
			return trackedTableCtrl.isCellDirty(row, cell);
		}

		function setCellDirty(cell, isDirty) {
			trackedTableCtrl.setCellDirty(row, cell, isDirty)
		}

		function setCellInvalid(cell, isInvalid) {
			trackedTableCtrl.setCellInvalid(row, cell, isInvalid)
		}
	}
})();

(function() {
	angular.module("myApp").directive("demoTrackedTableCell", demoTrackedTableCell);

	demoTrackedTableCell.$inject = [];

	function demoTrackedTableCell() {
		return {
			restrict: "A",
			priority: -1,
			scope: true,
			require: ["^demoTrackedTableRow", "ngForm"],
			controller: demoTrackedTableCellController
		};
	}

	demoTrackedTableCellController.$inject = ["$attrs", "$element", "$scope"];

	function demoTrackedTableCellController($attrs, $element, $scope) {
		var self = this;
		var cellFormCtrl = $element.controller("form");
		var cellName = cellFormCtrl.$name;
		var trackedTableRowCtrl = $element.controller("demoTrackedTableRow");

		if (trackedTableRowCtrl.isCellDirty(cellName)) {
			cellFormCtrl.$setDirty();
		} else {
			cellFormCtrl.$setPristine();
		}
		// note: we don't have to force setting validaty as angular will run validations
		// when we page back to a row that contains invalid data

		$scope.$watch(function() {
			return cellFormCtrl.$dirty;
		}, function(newValue, oldValue) {
			if (newValue === oldValue) return;

			trackedTableRowCtrl.setCellDirty(cellName, newValue);
		});

		$scope.$watch(function() {
			return cellFormCtrl.$invalid;
		}, function(newValue, oldValue) {
			if (newValue === oldValue) return;

			trackedTableRowCtrl.setCellInvalid(cellName, newValue);
		});
	}
})();

function  dataFactory() {
		return [{"id":1,"name":"Nissim","age":41,"money":454},{"id":2,"name":"Mariko","age":10,"money":-100},{"id":3,"name":"Mark","age":39,"money":291},{"id":4,"name":"Allen","age":85,"money":871},{"id":5,"name":"Dustin","age":10,"money":378},{"id":6,"name":"Macon","age":9,"money":128},{"id":7,"name":"Ezra","age":78,"money":11},{"id":8,"name":"Fiona","age":87,"money":285},{"id":9,"name":"Ira","age":7,"money":816},{"id":10,"name":"Barbara","age":46,"money":44},{"id":11,"name":"Lydia","age":56,"money":494},{"id":12,"name":"Carlos","age":80,"money":193}];
	}
