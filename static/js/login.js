$(document).ready(function(){
    M.AutoInit();
    Waves.displayEffect();

    $("#login-button").click(function(){
        M.toast({html: "Sending login request!"});

        $.ajax({
            url: "/