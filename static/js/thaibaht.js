
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