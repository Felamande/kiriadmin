window.onload = function() {
    Vue.use(VueResource)

    window.vm = new Vue({
        el:".container",
        data:{
            mdSrc:"",
            mdCvt:"",
            mdPrv:document.getElementById("md-preview")
        },
        components:{
            group:vuxGroup,
            cell:vuxCell,
            "x-textarea":vuxXTextarea,
            flexbox:vuxFlexbox,
            "flexbox-item":vuxFlexboxItem,
            "x-button":vuxXButton
        },
        methods:{
            postPreview:function() {
                this.$http.post("/editor/preview",{"content":this.mdSrc,"_xsrf":document.getElementsByName('_xsrf')[0].value},{headers:{"Content-Type":"application/json"}}).then(
                    function(response){
                        if(typeof response.data == 'object'){
                            this.handleError(response.data.err)
                            return
                        }
                        document.getElementById("md-preview").innerHTML = response.data
                    },function(error){}
                )
            },
            handleError:function(err){
                console.log(err)
            }
        }
    })
}

window.onbeforeunload = function(){
    return (function(){return "hello"})()
}