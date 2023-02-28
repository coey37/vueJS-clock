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
                    offer.remove();
                    M.toast({html: "Successfully cancelled trade."});
                } else {
                    M.toast({html: "Error cancelling trade, refresh the page."});
                }
            }
        });
    });

    $(".accept-button").click(function(){
        var offer = $(this).closest(".trade-offer");
        var id = parseInt(offer.attr("data-id"));

        M.toast({html: "Accepting trade offer."});

        $.ajax({
            url: "/panel/trade/accept",
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
                    offer.children(".status").text("Offer accepted");
                    M.toast({html: "Successfully accepted trade."});
                } else {
                    M.toast({html: "Error accepting trade, refresh the page."});
                }
            }
        });
    });
});
