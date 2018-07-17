var isEditable  = true;

var type;

function updateEditableCell(element) {
    var newValue = $(element).val()
    var table = $(element).closest("table").DataTable().table();
    var row = table.row($(element).parents('tr'));
    //console.log(row.data());
    //console.log(element.val());
    isEditable = true; 
    
    var isSend = true;
    var data = {
        hash_id: row.data().hash_id, 
        user_says: row.data().user_says,
        actual_intent: row.data().actual_intent,
        status: row.data().status,
        addition_to_df: row.data().addition_to_df,
    }
    console.log("new data : " + newValue)
    if(type == 6) {
        if(newValue === row.data().user_says){
            isSend = false;
        } else {
            data.user_says = newValue;
        }
    } else if(type == 7) {
        if( newValue === row.data().actual_intent)
            isSend = false;
        else
            data.actual_intent = newValue;
    } else if(type == 8) {
        if(newValue === row.data().status)
            isSend = false;
        else
            data.status = newValue;
    } else if(type == 9) {
        if(newValue === row.data().addition_to_df)
            isSend = false;
        else
            data.addition_to_df = newValue;
    }

    if(isSend){
        console.log("sending")
        console.log(data)
        sendData(data, table);
    } else{
        console.log("drawing")
        table.draw();
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
            $('#response pre').html( JSON.stringify( data ) );
            reloadDatatable(table)
        },
        error: function( jqXhr, textStatus, errorThrown ){
            console.log( errorThrown );
        }
    });
    console.log(data)
}

function reloadDatatable(datatable) {
    datatable.ajax.reload();
}
    
$(function () {
    var table = $("#user_table").DataTable({
        "info": true,
        "processing": true,
        "serverSide": true,
        "ordering": false,
        "ajax": "/form/data",
        "columns": [
            { data: 'date_timestamp' },
            { data: 'from_uid' },
            { data: 'intent_name' },
            { data: 'chat_from_user' },
            { data: 'score' },
            { data: 'chat_from_bot' },
            { data: 'user_says' },
            { data: 'actual_intent' },
            { data: 'status' },
            { data: 'addition_to_df' },
            { data: 'pic' },
            { data: 'hash_id'},
            { data: ''}
        ], 
        "columnDefs": [ { 
                "targets": -2, 
                "visible": false
            },
            {
                "targets": -1,
                "visible": false,
                "data": null,
                "defaultContent": "<button>Edit</button>"
        } ]
    });
/*
    $('#user_table tbody').on( 'click', 'button', function () {
        var data = table.row( $(this).parents('tr') ).data();
        alert( data.pic +"'s hashid is: "+ data.hash_id );
    } );*/


    $("#user_table tbody").on("click", "td", function() {
        var data = table.row( $(this).parents('tr') ).data();
        //alert( data.pic +"'s hashid is: "+ data.hash_id );

        var columnIndex = table.cell( this ).index().column;
        if(columnIndex > 5 && columnIndex < 10 && isEditable){
//            alert( 'Clicked on cell in visible column: '+columnIndex);
            var cell = table.cell(this).node();
            var oldData;
            if(!$(cell).find("input").length){
                if (columnIndex == 6) {
                    oldData = data.user_says;
                    type = 6;
                } else if(columnIndex == 7) {
                    oldData = data.actual_intent;
                    type = 7;
                } else if(columnIndex == 8) {
                    oldData = data.status;
                    type = 8;
                } else if(columnIndex == 9) {
                    oldData = data.addition_to_df;
                    type = 9;
                }
  //              alert("open input")

                
                var html = "<input class='editable' onblur='updateEditableCell($(this))' value='" + oldData + "'></input>"
                $(cell).html(html);
                $('#editable').focus();
                isEditable = false;
            }
            //reloadDatatable();
        }
    })

/*    var intervalId = setInterval(function(){ 
        $.ajax({
            url: "/visitor_count", 
            dataType: "json",
            success: handleResponse,
            error: handleErrorResponse
        }); 
    }, 2000);

    function handleResponse(data) {
        $("#visitor_count").text("Visitor Count : " + data.visitor_count)
    }

    function handleErrorResponse(request, status, error){
        console.log(error)
        console.log("visitor counter doesn't work!!!")
        clearInterval(intervalId)
    }*/
    
});

