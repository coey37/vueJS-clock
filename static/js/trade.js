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
                