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

function addCommasNonZeroPoint(nStr)
{
	nStr += '';
	x = nStr.split('.');
	x1 = x[0];
	x2 = x.length > 1 ? '.' + x[1] : '';
	var rgx = /(\d+)(\d{3})/;
	while (rgx.test(x1)) {
		x1 = x1.replace(rgx, '$1' + ',' + '$2');
    }     
	return x1 + x2;
}

function hideGlobalModal()
{
    $('#global-modal').modal("hide");
}
function showGlobalModal()
{
    $('#global-modal').modal("show");
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

var getBathText = function (inputNumber) {

    var getText = function (input) {
      var toNumber = input.toString();
      var numbers = toNumber.split('').reverse();

      var numberText = "/หนึ่ง/สอง/สาม/สี่/ห้า/หก/เจ็ด/แปด/เก้า/สิบ".split('/');
      var unitText = "/สิบ/ร้อย/พ้น/หมื่น/แสน/ล้าน".split('/');

      var output = "";
      for (var i = 0; i < numbers.length; i++) {
        var number = parseInt(numbers[i]);
        var text = numberText[number];
        var unit = unitText[i];

        if (number == 0)
          continue;

        if (i == 1 && number == 2) {
          output = "ยี่สิบ" + output;
          continue;
        }

        if (i == 1 && number == 1) {
          output = "สิบ" + output;
          continue;
        }


        output = text + unit + output;
      }

      return output;
    }

    var fullNumber = Math.floor(inputNumber);
    var decimal = inputNumber - fullNumber;

    if (decimal == 0) {

      return getText(fullNumber) + "บาทถ้วน";

    } else {

      // convert decimal into full number, need only 2 digits
      decimal = decimal * 100;
      decimal = Math.round(decimal);

      return getText(fullNumber) + "บาท" + getText(decimal) + "สตางค์";
    }

  };
