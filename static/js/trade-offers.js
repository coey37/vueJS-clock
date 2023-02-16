$(document).ready(function(){
    var classid = $(location).attr("href").split('/').pop();

    $(".cancel-button").click(function(){
        var offer = $(this).closest(".trade-offer");
        var id = parseInt(offer.attr("data-id"));

        M.toast({html: "Cancelling trade offer."});

        $.ajax({
            url: "/panel/trade/cancel",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                ID: id,
                CsrfSecret: CsrfSecret
            }),
            dataType: "json",
            success: function(r) {
                M.Toast.dismissAll(); // Clear all other toasts.
                if(r.success) {
         