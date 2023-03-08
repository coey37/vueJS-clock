$(document).ready(function(){
    var classid = $(location).attr("href").split('/').pop();

    $("#submit").click(function(){
        M.toast({html: "Sending trade offer."});

        $.ajax({
            url: "/panel/trade",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                User: $("#user").val(),
                Points: parseInt($("#points").val()),
                Classid: classid,
                CsrfSecret: CsrfSecret
            }),
            dataType: "json",
            success: function(r) {
                if(r.success) {
                    window.location.replace("/panel/trade-offers");
                } else {
                    M.Toast.dismissAll(); // Clear all other toasts.
                 