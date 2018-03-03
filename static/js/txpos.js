function addCommas(nStr)
{
	nStr += '';
	x = nStr.split('.');
	x1 = x[0];
	x2 = x.length > 1 ? '.' + x[1] : '';
	var rgx = /(\d+)(\d{3})/;
	while (rgx.test(x1)) {
		x1 = x1.replace(rgx, '$1' + ',' + '$2');
    }
    if(x2==''){
        x2 = ".00";
    }
	return x1 + x2;
}

function editNormal(id) {
    hideTopAlert();
    $.get("/normal/add/?entity="+ $("#entity").val() + "&id=" + id , function (data) {
        if (data.RetOK) 
        {
            showGlobalSmallModal();
            $('#small-global-modal-content').html(data.RetData);

        } else {
            showGlobalSmallModal();
            $('#small-global-modal-content').html(data.RetData);
        }
    });
}
 

function confirmDeleteGlobal(id,url) {
    hideTopAlert();
    hideGlobalDelete();
    $.get("/service/secure/json" , function (data) {
         $("#global-delete-xsrf").val(data)
    });
    $("#global-delete-id").val(id)
    $("#global-delete-url").val(url)
    $("#delete-global-modal").modal("show");
}
function deleteGlobal() {
    hideTopAlert();
    $.ajax({
        url: $("#global-delete-url").val() + "/" + $("#global-delete-id").val() ,
        type: 'DELETE',
        beforeSend: function (xhr) { xhr.setRequestHeader('X-Xsrftoken', $("#global-delete-xsrf").val()); },
        success: function (data) {
            if (data.RetOK) {
                showTopAlert(data.RetData, "success");
                $("#delete-global-modal").modal("hide");
                LoadData();
            } else {
                showGlobalDelete(data.RetData);
            }
        }
    });
}

function hideGlobalDelete(){
    $("#global-delete-alert").hide();
}

function showGlobalDelete(msg){
    $("#global-delete-alert").html(msg);
    $("#global-delete-alert").show();
}

function hideGlobalSmalModal()
{
    $('#small-global-modal').modal("hide");
}

function showGlobalSmallModal()
{
    $('#small-global-modal').modal("show");
}

function showTopAlert(alert,type)
{
    var html = `<div class="alert alert-`+type+` alert-dismissible"  >
    <button type="button" class="close" data-dismiss="alert" aria-hidden="true">×</button>
    `+ alert + `
    </div>`;
    $("#top-alert").html(html)
    $("#top-alert").fadeIn(500, function () {

    });
}

function hideTopAlert()
{
    $("#top-alert").fadeOut(500, function () {
    });
}

$(function () {
      $("#name-l").html($("#name-r").html().substring(0,20));
});
