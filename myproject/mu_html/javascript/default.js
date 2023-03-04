function onConfrm() { 
	var YorN = confirm("!您真的确定执行此操作吗!");
	if(YorN) return true;
	else return false;
}
function GotoSearchPage(mod){
	SearchTxt=document.getElementById(mod);
	window.location.href='/search/all/'+SearchTxt.value+"/0";
}
function DisplayContent(str){
	element=document.getElementById(str);
	if(element.style.display=="block"){
		element.style.display="none";
	}
	else{
		element.style.display="block";
	}
}
function StopEmptyForm(){
    if(document.getElementsByName("username")[0].value!="" && document.getElementsByName("password")[0].value!="")
        return true;
    else{
        alert("密码或账户为空")
        return false;
    }
}