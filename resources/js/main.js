var isEditable  = true;

var type;

function updateSelectCell(element) {
    isEditable = true;
    var newValue = $(element).val()
    var table = $(element).closest("table").DataTable().table();
    var row = table.row($(element).parents('tr'));
    var data = {
        hash_id: row.data().hash_id, 
        user_says: row.data().user_says,
        actual_intent: row.data().actual_intent,
        suggested_new_intent: row.data().suggested_new_intent
    }

    if(type === 7){
        if(newValue !== row.data().actual_intent){
            data.actual_intent = newValue
            sendData(data, table);
        } else {
            reloadDatatable(table)
        }
    } else if(type === 8){
        if(newValue !== row.data().suggested_new_intent){
            data.suggested_new_intent = newValue
            sendData(data, table);
        } else {
            reloadDatatable(table)
        }
    }
}

function updateEditableCell(element) {
    var newValue = $(element).val()
    var table = $(element).closest("table").DataTable().table();
    var row = table.row($(element).parents('tr'));
    isEditable = true; 
    if(row.data() === undefined){
        return
    }
    
    var isSend = true;
    var data = {
        hash_id: row.data().hash_id, 
        user_says: row.data().user_says,
        actual_intent: row.data().actual_intent,
        suggested_new_intent: row.data().suggested_new_intent
    }
    if(type == 6) {
        if(newValue === row.data().user_says){
            isSend = false;
        } else {
            data.user_says = newValue;
        }
    } 

    if(isSend){
        sendData(data, table);
    } else{
        reloadDatatable(table)
    }
}

function sendData(data, table) {
    $.ajax({
        url: '/form/ajax_post',
        dataType: 'json',
        type: 'post',
        contentType: 'application/json',
        data: JSON.stringify( data ),
        processData: false,
        success: function( data, textStatus, jQxhr ){
            reloadDatatable(table)
        },
        error: function( jqXhr, textStatus, errorThrown ){
            console.log( errorThrown );
        }
    });
}

function reloadDatatable(datatable) {
    datatable.ajax.reload(null, false);
    isEditable  = true;
    type = 0;

}
    
$(function () {
    var table = $("#user_table").DataTable({
        "info": true,
        "processing": true,
        "serverSide": true,
        "ordering": false,
        "ajax": "/form/data",
        "scrollX": true,
        "scrollY": 700,
        "drawCallback": function(settings) {
            isEditable = true; 
         },
        "columns": [
            { data: 'date_timestamp' },
            { data: 'from_uid' },
            { data: 'intent_name' },
            { data: 'chat_from_user' },
            { data: 'score' },
            { data: 'chat_from_bot' },
            { data: 'user_says' },
            { data: 'actual_intent' },
            { data: 'suggested_new_intent'},
            { data: 'status' },
            { data: 'addition_to_df' },
            { data: 'pic' },
            { data: 'hash_id'},
            { data: ''}
        ], 
        "columnDefs": [ { 
                "targets": -2, 
                "visible": false
            },{
                "targets": -1,
                "visible": false,
                "data": null,
                "defaultContent": "<button>Edit</button>"
            },{
                "className": "editable_field",
                "width": 400,
                "targets": [-6, -7, -8]
            }]
    });

    $("#user_table tbody").on("click", "td", function() {
        var data = table.row( $(this).parents('tr') ).data();
        console.log(data)
        var columnIndex = table.cell( this ).index().column;
        if(columnIndex > 5 && columnIndex < 9 && isEditable){
            var cell = table.cell(this).node();
            var oldData;
            if(!$(cell).find("input").length){
                var html;
                if (columnIndex == 6) {
                    oldData = data.user_says;
                    type = 6;
                    html = "<textarea class='editable' onblur='updateEditableCell($(this))' width: 100%; height: 100%>"+oldData+"</textarea>"
                    $(cell).html(html);
                    $('.editable').focus();
                } else if(columnIndex == 7) {
                    oldData = data.actual_intent;
                    type = 7;
                    html = "<select class='intent' style='width: 100%'><option value='"+oldData+"' selected='selected'>"+oldData+"</option></select>"
                    $(cell).html(html);
                    $('.intent').select2({
                        ajax: {
                            url: '/form/intent/ajax',
                            dataType: 'json'
                        },
                        width: 'resolve',
                        minimumInputLength: 3,
                        placeholder: "please choose.."
                    });

                    $(".intent").on("select2:close", function () { 
                        updateSelectCell($(this))
                     });
                     $('.intent').select2('open');
                } else if(columnIndex == 8) {
                    oldData = data.suggested_new_intent;
                    type = 8;

                    html = "<select class='suggest-intent' style='width: 100%'><option value='"+oldData+"' selected='selected'>"+oldData+"</option></select>"
                    $(cell).html(html);
                    $('.suggest-intent').select2({
                        ajax: {
                            url: '/form/suggest_intent/ajax',
                            dataType: 'json'
                        },
                        width: 'resolve',
                        minimumInputLength: 3,
                        placeholder: "please choose or input..",
                        tags: true
                    });

                    $(".suggest-intent").on("select2:close", function () { 
                        updateSelectCell($(this))
                     });
                     $('.suggest-intent').select2('open');
                } 
                
                isEditable = false;
            }
        }
    })
    
});

