$(document).ready(function(){
    M.AutoInit();
    Waves.displayEffect();

    $("#button").click(function(){
        $("#message").html("Sending forgot password message!");
 
        $.ajax({
            url: "/forgot-password",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                Email: $("#email").val(),
                Captcha: grecaptcha.getResponse()
            }),
            dataType: "json",
            success: function(r) {
                console.log(r.Code);
                switch(r.Code) {
                    case 0: {
                        $("#message").html("Successfully sent email to "