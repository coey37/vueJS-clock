$(document).ready(function(){
    var classid = $(location).attr("href").split('/').pop();

    $(".cancel-button").click(function(){
        var offer = $(this).closest(".trade-offer");
        var id = parseInt(offer.attr("data-id"));

        M.toast({html: