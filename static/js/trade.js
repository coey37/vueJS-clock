$(document).ready(function(){
    var classid = $(location).attr("href").split('/').pop();

    $("#submit").click(function(){
        M.toast({html: "Sending trade 