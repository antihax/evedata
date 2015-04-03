

	$(document).ready(function(){
 
		$("ul.subnav").parent().append("<span></span>"); //Only shows drop down trigger when js is enabled - Adds empty span tag after ul.subnav
	
		$("ul.topnav li span").click(function() { //When trigger is clicked...
		
			//Following events are applied to the subnav itself (moving subnav up and down)
			$(this).parent().find("ul.subnav").show(); //Drop down the subnav on click
 
			$(this).parent().hover(function() {
			}, function(){	
				$(this).parent().find("ul.subnav").slideUp('slow'); //When the mouse hovers out of the subnav, move it back up
			});
 
			//Following events are applied to the trigger (Hover events for the trigger)
			}).hover(function() { 
				$(this).addClass("subhover"); //On hover over, add class "subhover"
			}, function(){	//On Hover Out
				$(this).removeClass("subhover"); //On hover out, remove class "subhover"
		});
 	});

/************************************************************************************************************
	(C) www.dhtmlgoodies.com, November 2005
	
	This is a script from www.dhtmlgoodies.com. You will find this and a lot of other scripts at our website.	
	
	Terms of use:
	You are free to use this script as long as the copyright message is kept intact. However, you may not
	redistribute, sell or repost it without our permission.
	
	Thank you!
	
	www.dhtmlgoodies.com
	Alf Magne Kalleland
	
	************************************************************************************************************/	
	var arrayOfRolloverClasses = new Array();
	var arrayOfClickClasses = new Array();
	var activeRow = false;
	var activeRowClickArray = new Array();
	
	function highlightTableRow()
	{
		var tableObj = this.parentNode;
		if(tableObj.tagName!='TABLE')tableObj = tableObj.parentNode;

		if(this!=activeRow){
			this.setAttribute('origCl',this.className);
			this.origCl = this.className;
		}
		this.className = arrayOfRolloverClasses[tableObj.id];
		
		activeRow = this;
		
	}
	
	function clickOnTableRow()
	{
		var tableObj = this.parentNode;
		if(tableObj.tagName!='TABLE')tableObj = tableObj.parentNode;		
		
		if(activeRowClickArray[tableObj.id] && this!=activeRowClickArray[tableObj.id]){
			activeRowClickArray[tableObj.id].className='';
		}
		this.className = arrayOfClickClasses[tableObj.id];
		
		activeRowClickArray[tableObj.id] = this;
				
	}
	
	function resetRowStyle()
	{
		var tableObj = this.parentNode;
		if(tableObj.tagName!='TABLE')tableObj = tableObj.parentNode;

		if(activeRowClickArray[tableObj.id] && this==activeRowClickArray[tableObj.id]){
			this.className = arrayOfClickClasses[tableObj.id];
			return;	
		}
		
		var origCl = this.getAttribute('origCl');
		if(!origCl)origCl = this.origCl;
		this.className=origCl;
		
	}
		
	function addTableRolloverEffect(tableId,whichClass,whichClassOnClick)
	{
		arrayOfRolloverClasses[tableId] = whichClass;
		arrayOfClickClasses[tableId] = whichClassOnClick;
		
		var tableObj = document.getElementById(tableId);
		var tBody = tableObj.getElementsByTagName('TBODY');
		if(tBody){
			var rows = tBody[0].getElementsByTagName('TR');
		}else{
			var rows = tableObj.getElementsByTagName('TR');
		}
		for(var no=0;no<rows.length;no++){
			rows[no].onmouseover = highlightTableRow;
			rows[no].onmouseout = resetRowStyle;
			
			if(whichClassOnClick){
				rows[no].onclick = clickOnTableRow;	
			}
		}
		
	}

	var tableWidget_okToSort = true;
	var tableWidget_arraySort = new Array();
	tableWidget_tableCounter = 1;
	var activeColumn = new Array();
	var currentColumn = false;
	
	function sortNumeric(a,b){
		
		a = parseFloat(a.replace(/,/g,""));
		b = parseFloat(b.replace(/,/g,""));
		return a - b;
	}
	

	function sortString(a, b) {

		var aa = a.toLowerCase();
		var bb = b.toLowerCase();
		if (aa==bb) {
			return 0;
		} else if (aa<bb) {
			return -1;
		} else {
			return 1;
		}
	}
		
	function sortTable()
	{
		if(!tableWidget_okToSort)return;
		tableWidget_okToSort = false;
		/* Getting index of current column */
		var obj = this;
		var indexThis = 0;
		while(obj.previousSibling){
			obj = obj.previousSibling;
			if(obj.tagName=='TD')indexThis++;		
		}
		
		if(this.getAttribute('direction') || this.direction){
			direction = this.getAttribute('direction');
			if(navigator.userAgent.indexOf('Opera')>=0)direction = this.direction;
			if(direction=='ascending'){
				direction = 'descending';
				this.setAttribute('direction','descending');
				this.direction = 'descending';	
			}else{
				direction = 'ascending';
				this.setAttribute('direction','ascending');		
				this.direction = 'ascending';		
			}
		}else{
			direction = 'ascending';
			this.setAttribute('direction','ascending');
			this.direction = 'ascending';
		}
		
		var tableObj = this.parentNode.parentNode.parentNode;
		var tBody = tableObj.getElementsByTagName('TBODY')[0];
		
		var widgetIndex = tableObj.getAttribute('tableIndex');
		if(!widgetIndex)widgetIndex = tableObj.tableIndex;
		
		if(currentColumn)currentColumn.className='';
		document.getElementById('col' + widgetIndex + '_' + (indexThis+1)).className='highlightedColumn';
		currentColumn = document.getElementById('col' + widgetIndex + '_' + (indexThis+1));

			
		var sortMethod = tableWidget_arraySort[widgetIndex][indexThis]; // N = numeric, S = String
		if(activeColumn[widgetIndex] && activeColumn[widgetIndex]!=this){
			if(activeColumn[widgetIndex])activeColumn[widgetIndex].removeAttribute('direction');			
		}

		activeColumn[widgetIndex] = this;
		
		var cellArray = new Array();
		var cellObjArray = new Array();
		for(var no=1;no<tableObj.rows.length;no++){
			var content= tableObj.rows[no].cells[indexThis].innerHTML+'';
			cellArray.push(content);
			cellObjArray.push(tableObj.rows[no].cells[indexThis]);
		}
		
		if(sortMethod=='N'){
			cellArray = cellArray.sort(sortNumeric);
		}else{
			cellArray = cellArray.sort(sortString);
		}
		
		if(direction=='descending'){
			for(var no=cellArray.length;no>=0;no--){
				for(var no2=0;no2<cellObjArray.length;no2++){
					if(cellObjArray[no2].innerHTML == cellArray[no] && !cellObjArray[no2].getAttribute('allreadySorted')){
						cellObjArray[no2].setAttribute('allreadySorted','1');	
						tBody.appendChild(cellObjArray[no2].parentNode);				
					}				
				}			
			}
		}else{
			for(var no=0;no<cellArray.length;no++){
				for(var no2=0;no2<cellObjArray.length;no2++){
					if(cellObjArray[no2].innerHTML == cellArray[no] && !cellObjArray[no2].getAttribute('allreadySorted')){
						cellObjArray[no2].setAttribute('allreadySorted','1');	
						tBody.appendChild(cellObjArray[no2].parentNode);				
					}				
				}			
			}				
		}
		
		for(var no2=0;no2<cellObjArray.length;no2++){
			cellObjArray[no2].removeAttribute('allreadySorted');		
		}

		tableWidget_okToSort = true;
		
		
	}
	function initSortTable(objId,sortArray)
	{
		var obj = document.getElementById(objId);
		obj.setAttribute('tableIndex',tableWidget_tableCounter);
		obj.tableIndex = tableWidget_tableCounter;
		tableWidget_arraySort[tableWidget_tableCounter] = sortArray;
		var tHead = obj.getElementsByTagName('THEAD')[0];
		var cells = tHead.getElementsByTagName('TD');
		for(var no=0;no<cells.length;no++){
			if(sortArray[no]){
				cells[no].onclick = sortTable;	
			}else{
				cells[no].style.cursor = 'default';	
			}
		}		
		for(var no2=0;no2<sortArray.length;no2++){	/* Right align numeric cells */
			if(sortArray[no2] && sortArray[no2]=='N')obj.rows[0].cells[no2].style.textAlign='right';
		}		
		
		tableWidget_tableCounter++;
	}
		
