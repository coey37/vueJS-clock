$(document).ready(function(){
    M.AutoInit();
    Waves.displayEffect();

    $("#login-button").click(function(){
        M.toast({html: "Sending login request!"});

        $.ajax({
            url: "/login",
            type: "POST",
            contentType: "application/json; charset=utf-8",
            data: JSON.stringify({
                Email: $("#email").val(),
                Password: $("#password").val(),
                Captcha: grecaptcha.getResponse()
            }),
            dataType: "json",
            success: function(r) {
