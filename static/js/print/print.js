 var pageHeight = 1400;
        $(document).ready(function () {
            $(document).on('click', '.keep-open', function (e) {
                e.stopPropagation();
            });
            printPreview();
            resize_to_fit("headerarea0", 30);  
        });

        function printPreview() {
            var pageIndex = 1,
                pageIndexSrting,
                pageIndexClass,
                currentPage;
            var count = 0;
            var countpage = 0;
            for (var a = 0; a < 1; a++) {

                var pagedoc = '.print-page-' + a;
                var arrItems = [];
                var rowItems = $(pagedoc + ' .print-table-item'),
                    rowHeader = $(pagedoc + ' .print-items-header'),
                    rowFooter = $(pagedoc + ' .print-summary');
                var printPreview = $('.print-preview');
                var headerSection = $(pagedoc + ' .print-header-section'),
                    mainSection = $(pagedoc + ' .print-main-section'),
                    footerSection = $(pagedoc + ' .print-footer-section'),
                    tableRow = $(pagedoc + ' .table-items');
                var header, main, footer;
                var main = "<div class='print-main-section'>";
                main += "<div class='print-items'>";
                main += "<table class='table-items'><tbody></tbody></table></div></div>";


                var headerHeight;
                var infoHeight;
                var contentHeight; // Table item
                var footerHeight;

                headerHeight = Math.round($(pagedoc + ' .print-header').outerHeight(true));

                infoHeight = Math.round($(pagedoc + ' .print-info-section').outerHeight(true));
                footerHeight = Math.round($(pagedoc + ' .print-footer-section').outerHeight(true));
                contentHeight = pageHeight - (headerHeight + infoHeight + footerHeight);
                var columnTable;
                var contentTable = contentHeight;
                var footerTable;
                var egdeTable;
                columnTable = $('#items-column-' + a).outerHeight(true)
                footerTable = $('#summary-subtotal-' + a).outerHeight(true) + $('#summary-before-tax-' + a).outerHeight(true);
                footerTable += $('#summary-tax-' + a).outerHeight(true) + $('#summary-total-' + a).outerHeight(true) + $('#summary-shippingamount-' + a).outerHeight(true);
                contentTable = contentTable - (columnTable + footerTable);
                egdeTable = contentTable + footerTable;
                var objRow = $(pagedoc + ' .print-table-item').clone();
                var rowItems = $(pagedoc + ' .print-table-item');
                var arrayPages = splitPage(objRow);
                createPage(arrayPages);
            }
            function splitPage(rows) {
                var pages = [];
                var rowCountHeight = 0;
                var flagEgde = 0;

                var row_f = rowFooter.clone();
                pages.push('page');
                var emptyheight = 40;

                var rowTemp = rowItems.clone();
                for (i = 0; i < rowItems.length; i++) {

                    rowCountHeight += rowItems[i].clientHeight;

                    if (rowCountHeight < contentTable) {
                        pages.push(rowTemp[i]);
                        if (i == (rowItems.length - 1) && rowCountHeight + emptyheight < contentTable) {
                            while (rowCountHeight + emptyheight < contentTable) {
                                var emptyrow = document.createElement("tr");
                                {
                                    var y = document.createElement("td");
                                    y.className = "item-no";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }
                                {
                                    var y = document.createElement("td");
                                    y.className = "item-id productcodearea";
                                    y.style.display = "";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }
                                {
                                    var y = document.createElement("td");
                                    y.className = "item-name";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }
                                {
                                    var y = document.createElement("td");
                                    y.className = "item-value";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }
                                {
                                    var y = document.createElement("td");
                                    y.className = "item-cost";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }

                                {
                                    var y = document.createElement("td");
                                    y.className = "item-sale";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }

                                {
                                    var y = document.createElement("td");
                                    y.className = "item-total";
                                    var t = document.createTextNode("\u00A0");
                                    y.appendChild(t);
                                    emptyrow.appendChild(y);
                                }
                                rowCountHeight += emptyheight;
                                pages.push(emptyrow);
                            }
                        }
                    } else if ((rowCountHeight > contentTable) && (rowCountHeight < egdeTable)) {
                        flagEgde++;
                        if (flagEgde == 1) {
                            pages.push('egde');
                            pages.push(rowTemp[i]);
                        } else {
                            pages.push(rowTemp[i]);
                        }

                        // If last item and too height to print slpit to new page
                        if (i == (rowItems.length - 1)) {
                            // Replace page to egde
                            var indexEgde = pages.lastIndexOf('egde');
                            pages[indexEgde] = 'page';
                        }

                        // New Page
                    } else {
                        flagEgde = 0;

                        // if have egde it mean is a long page
                        // remove egde fromt arry of pages
                        var sliceIndex = pages.indexOf('egde');
                        var arrayPrev = pages.slice(0, sliceIndex - 1);
                        var arrayNext = pages.slice(sliceIndex + 1, pages.length);
                        var arrayTemp = arrayPrev.concat(arrayNext);

                        pages = arrayTemp;

                        pages.push('page');
                        pages.push(rowTemp[i]);
                        // Reset
                        rowCountHeight = 0;
                    }
                    pages.push(row_f);
                }
                return pages;
            }
            function createPage(pages) {
                var startindex = 0;
                if (count != 0) {
                    startindex = countpage;
                }
                var allpage = 0;
                for (var c = startindex; c < pages.length; c++) {
                    if (pages[c] == 'page') {
                        allpage++;
                    }
                }
                var cpage = 0;
                for (index = startindex; index < pages.length; index++) {

                    // Initial Page
                    if (pages[index] == 'page') {

                        pageIndexSrting = "print-preview-page-" + pageIndex;
                        printPreview.append("<div class='" + pageIndexSrting + "'></div> ");
                        pageIndex++;

                        // Get Index of Page
                        currentPage = $('.' + pageIndexSrting);

                        header = headerSection.clone();
                        row = rowHeader.clone();

                        cpage++;
                        //currentPage.append("<span class='pull-right'>" + "หน้า" + " " + cpage + "/" + allpage + "</span>");
                        currentPage.append(header);
                        currentPage.append(main);
                        $('.' + pageIndexSrting + ' .table-items tbody').append(row);

                        if (pageIndex != 1) {
                            footer = footerSection.clone();
                            currentPage.append(footer);
                        }

                    } else {
                        $('.' + pageIndexSrting + ' .table-items tbody').append(pages[index]);
                    }
                }
                countpage += pages.length;
            }
        }
        function resize_to_fit(arr, countfont) {
            var jEl = $('#' + arr + ' span');
            var fontsize = jEl.css('font-size');
            if (fontsize <= 1)
                return;
            if (jEl.width() > $('#' + arr).width()) {
                jEl.css('fontSize', parseFloat(fontsize) - 1);
                if (countfont > 0)
                    resize_to_fit(arr, countfont - 1);
            }
        }