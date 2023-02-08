$(document).ready(function(){
    var classid = $(location).attr("href").split('/').pop();

    $(".cancel-button").click(function(){
        var offer = $(this).closest(".