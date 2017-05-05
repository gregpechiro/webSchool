package web

import "html/template"

var UAJS template.JS = `!function(i,s){"use strict";var e="0.7.10",o="",r="?",n="function",a="undefined",t="object",w="string",l="major",d="model",p="name",c="type",u="vendor",m="version",f="architecture",b="console",g="mobile",h="tablet",v="smarttv",x="wearable",y="embedded",k={extend:function(i,s){var e={};for(var o in i)s[o]&&s[o].length%2===0?e[o]=s[o].concat(i[o]):e[o]=i[o];return e},has:function(i,s){return"string"==typeof i?-1!==s.toLowerCase().indexOf(i.toLowerCase()):!1},lowerize:function(i){return i.toLowerCase()},major:function(i){return typeof i===w?i.split(".")[0]:s}},A={rgx:function(){for(var i,e,o,r,w,l,d,p=0,c=arguments;p<c.length&&!l;){var u=c[p],m=c[p+1];if(typeof i===a){i={};for(r in m)m.hasOwnProperty(r)&&(w=m[r],typeof w===t?i[w[0]]=s:i[w]=s)}for(e=o=0;e<u.length&&!l;)if(l=u[e++].exec(this.getUA()))for(r=0;r<m.length;r++)d=l[++o],w=m[r],typeof w===t&&w.length>0?2==w.length?typeof w[1]==n?i[w[0]]=w[1].call(this,d):i[w[0]]=w[1]:3==w.length?typeof w[1]!==n||w[1].exec&&w[1].test?i[w[0]]=d?d.replace(w[1],w[2]):s:i[w[0]]=d?w[1].call(this,d,w[2]):s:4==w.length&&(i[w[0]]=d?w[3].call(this,d.replace(w[1],w[2])):s):i[w]=d?d:s;p+=2}return i},str:function(i,e){for(var o in e)if(typeof e[o]===t&&e[o].length>0){for(var n=0;n<e[o].length;n++)if(k.has(e[o][n],i))return o===r?s:o}else if(k.has(e[o],i))return o===r?s:o;return i}},E={browser:{oldsafari:{version:{"1.0":"/8",1.2:"/1",1.3:"/3","2.0":"/412","2.0.2":"/416","2.0.3":"/417","2.0.4":"/419","?":"/"}}},device:{amazon:{model:{"Fire Phone":["SD","KF"]}},sprint:{model:{"Evo Shift 4G":"7373KT"},vendor:{HTC:"APA",Sprint:"Sprint"}}},os:{windows:{version:{ME:"4.90","NT 3.11":"NT3.51","NT 4.0":"NT4.0",2000:"NT 5.0",XP:["NT 5.1","NT 5.2"],Vista:"NT 6.0",7:"NT 6.1",8:"NT 6.2",8.1:"NT 6.3",10:["NT 6.4","NT 10.0"],RT:"ARM"}}}},S={browser:[[/(opera\smini)\/([\w\.-]+)/i,/(opera\s[mobiletab]+).+version\/([\w\.-]+)/i,/(opera).+version\/([\w\.]+)/i,/(opera)[\/\s]+([\w\.]+)/i],[p,m],[/(OPiOS)[\/\s]+([\w\.]+)/i],[[p,"Opera Mini"],m],[/\s(opr)\/([\w\.]+)/i],[[p,"Opera"],m],[/(kindle)\/([\w\.]+)/i,/(lunascape|maxthon|netfront|jasmine|blazer)[\/\s]?([\w\.]+)*/i,/(avant\s|iemobile|slim|baidu)(?:browser)?[\/\s]?([\w\.]*)/i,/(?:ms|\()(ie)\s([\w\.]+)/i,/(rekonq)\/([\w\.]+)*/i,/(chromium|flock|rockmelt|midori|epiphany|silk|skyfire|ovibrowser|bolt|iron|vivaldi|iridium|phantomjs)\/([\w\.-]+)/i],[p,m],[/(trident).+rv[:\s]([\w\.]+).+like\sgecko/i],[[p,"IE"],m],[/(edge)\/((\d+)?[\w\.]+)/i],[p,m],[/(yabrowser)\/([\w\.]+)/i],[[p,"Yandex"],m],[/(comodo_dragon)\/([\w\.]+)/i],[[p,/_/g," "],m],[/(chrome|omniweb|arora|[tizenoka]{5}\s?browser)\/v?([\w\.]+)/i,/(qqbrowser)[\/\s]?([\w\.]+)/i],[p,m],[/(uc\s?browser)[\/\s]?([\w\.]+)/i,/ucweb.+(ucbrowser)[\/\s]?([\w\.]+)/i,/JUC.+(ucweb)[\/\s]?([\w\.]+)/i],[[p,"UCBrowser"],m],[/(dolfin)\/([\w\.]+)/i],[[p,"Dolphin"],m],[/((?:android.+)crmo|crios)\/([\w\.]+)/i],[[p,"Chrome"],m],[/XiaoMi\/MiuiBrowser\/([\w\.]+)/i],[m,[p,"MIUI Browser"]],[/android.+version\/([\w\.]+)\s+(?:mobile\s?safari|safari)/i],[m,[p,"Android Browser"]],[/FBAV\/([\w\.]+);/i],[m,[p,"Facebook"]],[/fxios\/([\w\.-]+)/i],[m,[p,"Firefox"]],[/version\/([\w\.]+).+?mobile\/\w+\s(safari)/i],[m,[p,"Mobile Safari"]],[/version\/([\w\.]+).+?(mobile\s?safari|safari)/i],[m,p],[/webkit.+?(mobile\s?safari|safari)(\/[\w\.]+)/i],[p,[m,A.str,E.browser.oldsafari.version]],[/(konqueror)\/([\w\.]+)/i,/(webkit|khtml)\/([\w\.]+)/i],[p,m],[/(navigator|netscape)\/([\w\.-]+)/i],[[p,"Netscape"],m],[/(swiftfox)/i,/(icedragon|iceweasel|camino|chimera|fennec|maemo\sbrowser|minimo|conkeror)[\/\s]?([\w\.\+]+)/i,/(firefox|seamonkey|k-meleon|icecat|iceape|firebird|phoenix)\/([\w\.-]+)/i,/(mozilla)\/([\w\.]+).+rv\:.+gecko\/\d+/i,/(polaris|lynx|dillo|icab|doris|amaya|w3m|netsurf|sleipnir)[\/\s]?([\w\.]+)/i,/(links)\s\(([\w\.]+)/i,/(gobrowser)\/?([\w\.]+)*/i,/(ice\s?browser)\/v?([\w\._]+)/i,/(mosaic)[\/\s]([\w\.]+)/i],[p,m]],cpu:[[/(?:(amd|x(?:(?:86|64)[_-])?|wow|win)64)[;\)]/i],[[f,"amd64"]],[/(ia32(?=;))/i],[[f,k.lowerize]],[/((?:i[346]|x)86)[;\)]/i],[[f,"ia32"]],[/windows\s(ce|mobile);\sppc;/i],[[f,"arm"]],[/((?:ppc|powerpc)(?:64)?)(?:\smac|;|\))/i],[[f,/ower/,"",k.lowerize]],[/(sun4\w)[;\)]/i],[[f,"sparc"]],[/((?:avr32|ia64(?=;))|68k(?=\))|arm(?:64|(?=v\d+;))|(?=atmel\s)avr|(?:irix|mips|sparc)(?:64)?(?=;)|pa-risc)/i],[[f,k.lowerize]]],device:[[/\((ipad|playbook);[\w\s\);-]+(rim|apple)/i],[d,u,[c,h]],[/applecoremedia\/[\w\.]+ \((ipad)/],[d,[u,"Apple"],[c,h]],[/(apple\s{0,1}tv)/i],[[d,"Apple TV"],[u,"Apple"]],[/(archos)\s(gamepad2?)/i,/(hp).+(touchpad)/i,/(kindle)\/([\w\.]+)/i,/\s(nook)[\w\s]+build\/(\w+)/i,/(dell)\s(strea[kpr\s\d]*[\dko])/i],[u,d,[c,h]],[/(kf[A-z]+)\sbuild\/[\w\.]+.*silk\//i],[d,[u,"Amazon"],[c,h]],[/(sd|kf)[0349hijorstuw]+\sbuild\/[\w\.]+.*silk\//i],[[d,A.str,E.device.amazon.model],[u,"Amazon"],[c,g]],[/\((ip[honed|\s\w*]+);.+(apple)/i],[d,u,[c,g]],[/\((ip[honed|\s\w*]+);/i],[d,[u,"Apple"],[c,g]],[/(blackberry)[\s-]?(\w+)/i,/(blackberry|benq|palm(?=\-)|sonyericsson|acer|asus|dell|huawei|meizu|motorola|polytron)[\s_-]?([\w-]+)*/i,/(hp)\s([\w\s]+\w)/i,/(asus)-?(\w+)/i],[u,d,[c,g]],[/\(bb10;\s(\w+)/i],[d,[u,"BlackBerry"],[c,g]],[/android.+(transfo[prime\s]{4,10}\s\w+|eeepc|slider\s\w+|nexus 7)/i],[d,[u,"Asus"],[c,h]],[/(sony)\s(tablet\s[ps])\sbuild\//i,/(sony)?(?:sgp.+)\sbuild\//i],[[u,"Sony"],[d,"Xperia Tablet"],[c,h]],[/(?:sony)?(?:(?:(?:c|d)\d{4})|(?:so[-l].+))\sbuild\//i],[[u,"Sony"],[d,"Xperia Phone"],[c,g]],[/\s(ouya)\s/i,/(nintendo)\s([wids3u]+)/i],[u,d,[c,b]],[/android.+;\s(shield)\sbuild/i],[d,[u,"Nvidia"],[c,b]],[/(playstation\s[34portablevi]+)/i],[d,[u,"Sony"],[c,b]],[/(sprint\s(\w+))/i],[[u,A.str,E.device.sprint.vendor],[d,A.str,E.device.sprint.model],[c,g]],[/(lenovo)\s?(S(?:5000|6000)+(?:[-][\w+]))/i],[u,d,[c,h]],[/(htc)[;_\s-]+([\w\s]+(?=\))|\w+)*/i,/(zte)-(\w+)*/i,/(alcatel|geeksphone|huawei|lenovo|nexian|panasonic|(?=;\s)sony)[_\s-]?([\w-]+)*/i],[u,[d,/_/g," "],[c,g]],[/(nexus\s9)/i],[d,[u,"HTC"],[c,h]],[/[\s\(;](xbox(?:\sone)?)[\s\);]/i],[d,[u,"Microsoft"],[c,b]],[/(kin\.[onetw]{3})/i],[[d,/\./g," "],[u,"Microsoft"],[c,g]],[/\s(milestone|droid(?:[2-4x]|\s(?:bionic|x2|pro|razr))?(:?\s4g)?)[\w\s]+build\//i,/mot[\s-]?(\w+)*/i,/(XT\d{3,4}) build\//i,/(nexus\s[6])/i],[d,[u,"Motorola"],[c,g]],[/android.+\s(mz60\d|xoom[\s2]{0,2})\sbuild\//i],[d,[u,"Motorola"],[c,h]],[/android.+((sch-i[89]0\d|shw-m380s|gt-p\d{4}|gt-n8000|sgh-t8[56]9|nexus 10))/i,/((SM-T\w+))/i],[[u,"Samsung"],d,[c,h]],[/((s[cgp]h-\w+|gt-\w+|galaxy\snexus|sm-n900))/i,/(sam[sung]*)[\s-]*(\w+-?[\w-]*)*/i,/sec-((sgh\w+))/i],[[u,"Samsung"],d,[c,g]],[/(samsung);smarttv/i],[u,d,[c,v]],[/\(dtv[\);].+(aquos)/i],[d,[u,"Sharp"],[c,v]],[/sie-(\w+)*/i],[d,[u,"Siemens"],[c,g]],[/(maemo|nokia).*(n900|lumia\s\d+)/i,/(nokia)[\s_-]?([\w-]+)*/i],[[u,"Nokia"],d,[c,g]],[/android\s3\.[\s\w;-]{10}(a\d{3})/i],[d,[u,"Acer"],[c,h]],[/android\s3\.[\s\w;-]{10}(lg?)-([06cv9]{3,4})/i],[[u,"LG"],d,[c,h]],[/(lg) netcast\.tv/i],[u,d,[c,v]],[/(nexus\s[45])/i,/lg[e;\s\/-]+(\w+)*/i],[d,[u,"LG"],[c,g]],[/android.+(ideatab[a-z0-9\-\s]+)/i],[d,[u,"Lenovo"],[c,h]],[/linux;.+((jolla));/i],[u,d,[c,g]],[/((pebble))app\/[\d\.]+\s/i],[u,d,[c,x]],[/android.+;\s(glass)\s\d/i],[d,[u,"Google"],[c,x]],[/android.+(\w+)\s+build\/hm\1/i,/android.+(hm[\s\-_]*note?[\s_]*(?:\d\w)?)\s+build/i,/android.+(mi[\s\-_]*(?:one|one[\s_]plus)?[\s_]*(?:\d\w)?)\s+build/i],[[d,/_/g," "],[u,"Xiaomi"],[c,g]],[/\s(tablet)[;\/\s]/i,/\s(mobile)[;\/\s]/i],[[c,k.lowerize],u,d]],engine:[[/windows.+\sedge\/([\w\.]+)/i],[m,[p,"EdgeHTML"]],[/(presto)\/([\w\.]+)/i,/(webkit|trident|netfront|netsurf|amaya|lynx|w3m)\/([\w\.]+)/i,/(khtml|tasman|links)[\/\s]\(?([\w\.]+)/i,/(icab)[\/\s]([23]\.[\d\.]+)/i],[p,m],[/rv\:([\w\.]+).*(gecko)/i],[m,p]],os:[[/microsoft\s(windows)\s(vista|xp)/i],[p,m],[/(windows)\snt\s6\.2;\s(arm)/i,/(windows\sphone(?:\sos)*|windows\smobile|windows)[\s\/]?([ntce\d\.\s]+\w)/i],[p,[m,A.str,E.os.windows.version]],[/(win(?=3|9|n)|win\s9x\s)([nt\d\.]+)/i],[[p,"Windows"],[m,A.str,E.os.windows.version]],[/\((bb)(10);/i],[[p,"BlackBerry"],m],[/(blackberry)\w*\/?([\w\.]+)*/i,/(tizen)[\/\s]([\w\.]+)/i,/(android|webos|palm\sos|qnx|bada|rim\stablet\sos|meego|contiki)[\/\s-]?([\w\.]+)*/i,/linux;.+(sailfish);/i],[p,m],[/(symbian\s?os|symbos|s60(?=;))[\/\s-]?([\w\.]+)*/i],[[p,"Symbian"],m],[/\((series40);/i],[p],[/mozilla.+\(mobile;.+gecko.+firefox/i],[[p,"Firefox OS"],m],[/(nintendo|playstation)\s([wids34portablevu]+)/i,/(mint)[\/\s\(]?(\w+)*/i,/(mageia|vectorlinux)[;\s]/i,/(joli|[kxln]?ubuntu|debian|[open]*suse|gentoo|(?=\s)arch|slackware|fedora|mandriva|centos|pclinuxos|redhat|zenwalk|linpus)[\/\s-]?([\w\.-]+)*/i,/(hurd|linux)\s?([\w\.]+)*/i,/(gnu)\s?([\w\.]+)*/i],[p,m],[/(cros)\s[\w]+\s([\w\.]+\w)/i],[[p,"Chromium OS"],m],[/(sunos)\s?([\w\.]+\d)*/i],[[p,"Solaris"],m],[/\s([frentopc-]{0,4}bsd|dragonfly)\s?([\w\.]+)*/i],[p,m],[/(ip[honead]+)(?:.*os\s([\w]+)*\slike\smac|;\sopera)/i],[[p,"iOS"],[m,/_/g,"."]],[/(mac\sos\sx)\s?([\w\s\.]+\w)*/i,/(macintosh|mac(?=_powerpc)\s)/i],[[p,"Mac OS"],[m,/_/g,"."]],[/((?:open)?solaris)[\/\s-]?([\w\.]+)*/i,/(haiku)\s(\w+)/i,/(aix)\s((\d)(?=\.|\)|\s)[\w\.]*)*/i,/(plan\s9|minix|beos|os\/2|amigaos|morphos|risc\sos|openvms)/i,/(unix)\s?([\w\.]+)*/i],[p,m]]},T=function(s,e){if(!(this instanceof T))return new T(s,e).getResult();var r=s||(i&&i.navigator&&i.navigator.userAgent?i.navigator.userAgent:o),n=e?k.extend(S,e):S;return this.getBrowser=function(){var i=A.rgx.apply(this,n.browser);return i.major=k.major(i.version),i},this.getCPU=function(){return A.rgx.apply(this,n.cpu)},this.getDevice=function(){return A.rgx.apply(this,n.device)},this.getEngine=function(){return A.rgx.apply(this,n.engine)},this.getOS=function(){return A.rgx.apply(this,n.os)},this.getResult=function(){return{ua:this.getUA(),browser:this.getBrowser(),engine:this.getEngine(),os:this.getOS(),device:this.getDevice(),cpu:this.getCPU()}},this.getUA=function(){return r},this.setUA=function(i){return r=i,this},this};T.VERSION=e,T.BROWSER={NAME:p,MAJOR:l,VERSION:m},T.CPU={ARCHITECTURE:f},T.DEVICE={MODEL:d,VENDOR:u,TYPE:c,CONSOLE:b,MOBILE:g,SMARTTV:v,TABLET:h,WEARABLE:x,EMBEDDED:y},T.ENGINE={NAME:p,VERSION:m},T.OS={NAME:p,VERSION:m},typeof exports!==a?(typeof module!==a&&module.exports&&(exports=module.exports=T),exports.UAParser=T):typeof define===n&&define.amd?define("ua-parser-js",[],function(){return T}):i.UAParser=T;var N=i.jQuery||i.Zepto;if(typeof N!==a){var O=new T;N.ua=O.getResult(),N.ua.get=function(){return O.getUA()},N.ua.set=function(i){O.setUA(i);var s=O.getResult();for(var e in s)N.ua[e]=s[e]}}}("object"==typeof window?window:this);`

var activityPage = `<!DOCTYPE html>
<html>
    <head>
        <meta charset="utf-8"/>
        <meta http-equiv="X-UA-Compatible" content="IE=edge">
        <meta name="viewport" content="width=device-width, initial-scale=1.0, user-scalable=no"/>
        <meta http-equiv="Content-Type" content="text/html; charset=utf-8" />
        <!--[if lt IE 9]>
        <script src="//html5shim.googlecode.com/svn/trunk/html5.js"></script>
        <![endif]-->
        <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/bootstrap/3.3.5/css/bootstrap.min.css">
        <link rel="stylesheet" href="//maxcdn.bootstrapcdn.com/font-awesome/4.3.0/css/font-awesome.min.css">
        <link rel="stylesheet" href="//cdn.datatables.net/1.10.12/css/dataTables.bootstrap.min.css" >
        <style media="screen">
            tr#search input {
                width: 100%;
                padding: 3px;
                box-sizing: border-box;
            }

            tr#search th {
                padding-right: 8px;
            }

            .clickable {
                cursor: pointer;
            }
        </style>
        <script src="//ajax.googleapis.com/ajax/libs/jquery/2.1.1/jquery.min.js"></script>
        <script type="text/javascript">
            {{ .uaJs }}
        </script>
    </head>
    <body>
        <div class="container-fluid" style="margin-top:10px;">
            <div class="row">
                <div class="col-lg-12">
                    <table id="activities" class="table table-bordered">
                        <thead>
                            <tr id="search">
                                <th></th>
                                <th>ipAddress</th>
                                <th>userId</th>
                                <th></th>
                                <th>location</th>
                                <th>lastAction</th>
                                <th>browser</th>
                                <th>engine</th>
                                <th>os</th>
                                <th>cpu</th>
                            </tr>
                            <tr>
                                <th>Time Stamp</th>
                                <th>IP address</th>
                                <th>User Id</th>
                                <th>User Role</th>
                                <th>Location</th>
                                <th>Last Action</th>
                                <th>Browser</th>
                                <th>Engine</th>
                                <th>OS</th>
                                <th>CPU</th>
                            </tr>
                        </thead>
                        <tbody>
                            {{ range $i, $activity := .activities }}
                                {{ if ne $activity.UserAgent "Auvik HTTP Monitor" }}
                                    <tr>
                                        <td>{{ $activity.Time }}</td>
                                        <td>{{ $activity.Ip }}</td>
                                        <td>{{ $activity.Session.GetId }}</td>
                                        <td>{{ $activity.Session.GetRole }}</td>
                                        <td id="loc{{ $i }}"></td>
                                        <td>
                                            {{ $activity.Actions.Recent.Method }} {{ $activity.Actions.Recent.Path }} {{ $activity.Actions.Recent.Time }}
                                            <button data-actions="{{ $activity.Actions.ToJson }}" class="more btn btn-primary btn-xs pull-right">More</button>
                                        </td>
                                        <td id="browser{{ $i }}"></td>
                                        <td id="engine{{ $i }}"></td>
                                        <td id="os{{ $i }}"></td>
                                        <td id="cpu{{ $i }}"></td>
                                        <script type="text/javascript">
                                            if ('{{ $activity.UserAgent }}' !== 'Auvik HTTP Monitor') {
                                                var p = new UAParser('{{ $activity.UserAgent }}');
                                                $('td#browser{{ $i }}').text(p.getBrowser().name + ' ' + p.getBrowser().version);
                                                $('td#engine{{ $i }}').text(p.getEngine().name + ' ' + p.getEngine().version);
                                                $('td#os{{ $i }}').text(p.getOS().name + ' ' + p.getOS().version);
                                                $('td#cpu{{ $i }}').text(p.getCPU().architecture);
                                            }

                                            $.getJSON('http://ipinfo.io/{{ $activity.Ip }}', function(data){
                                                $('td#loc{{ $i }}').text(((data.city !== '' && data.city !== undefined ) ? data.city : '') + ((data.region !== '' && data.region !== undefined ) ? (' ' + data.region) : ''));
                                            })
                                        </script>
                                    </tr>
                                {{ end }}
                            {{ end }}
                        </tbody>
                    </table>
                </div>
            </div>
        </div>

        <div class="modal fade" id="actionModal" tabindex="-1" role="dialog" aria-labelledby="actionModalLabel">
            <div class="modal-dialog" role="document">
                <div class="modal-content">
                    <div class="modal-header">
                        <button type="button" class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></button>
                        <h4 class="modal-title" id="myModalLabel">Actions</h4>
                    </div>
                    <div class="modal-body" id="actionModalBody">
                        <table class="table">
                            <tbody id="actionModalTable">

                            </tbody>
                        </table>
                    </div>
                </div>
            </div>
        </div>

        <script src="//maxcdn.bootstrapcdn.com/bootstrap/3.2.0/js/bootstrap.min.js"></script>
        <script src="//cdn.datatables.net/1.10.12/js/jquery.dataTables.min.js" charset="utf-8"></script>
        <script src="//cdn.datatables.net/1.10.12/js/dataTables.bootstrap.min.js" charset="utf-8"></script>
        <script type="text/javascript">
            var table = $('#activities').DataTable({
                "lengthMenu":[10,15,20],
                "order": [[ 0, "desc" ]],
                "columnDefs": [
                    { "orderable": false,        "targets": [2] },
                    { "name"     : "ipAddress",  "targets": 1 },
                    { "name"     : "userId",     "targets": 2 },
                    { "name"     : "location",   "targets": 4 },
                    { "name"     : "lastAction", "targets": 5 },
                    { "name"     : "browser",    "targets": 6 },
                    { "name"     : "engine",     "targets": 7 },
                    { "name"     : "os",         "targets": 8 },
                    { "name"     : "cpu",        "targets": 9 }
                ]
            });

            $('button.more').click(function() {
                actions = JSON.parse($(this).attr('data-actions'));
                var ac = '';
                for (var i = 0; i < actions.length; i++) {
                    ac += '<tr>' +
                            '<td>' + actions[i].method + '</td>' +
                            '<td>' + actions[i].path + '</td>' +
                            '<td>' + actions[i].time + '</td>' +
                        '</tr>'
                }
                $('#actionModalTable').html(ac);
                $('#actionModal').modal('show');
            });

            $(document).on('click', 'tr.clickable', function() {
                if (this.getAttribute('data-target') === '_blank') {
                    window.open(this.getAttribute('data-url'));
                    return
                }
                window.location.href = this.getAttribute('data-url');
            });

            $('#filter').on( 'keyup', function () {
                if (table !== undefined && !$.isEmptyObject(table)) {
                    table.search( this.value ).draw();
                }
            });

            $('tr#search th').each(function () {
                var title = $(this).text();
                if (title != '') {
                    $(this).html('<input id="columnSearch" data-column="' + title + '" type="text"/>');
                }
            });

            $('input#columnSearch').keyup(function() {
                var name = $(this).attr('data-column');
                var column = table.column(name + ':name');
                if (column.search() !== this.value ) {
                    column.search(this.value).draw();
                }
            });

        </script>

    </body>
</html>
`
