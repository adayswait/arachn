import Vue from 'vue';
import Router from 'vue-router';
import Home from "@/pages/Home.vue"
import DevOps from "@/pages/DevOps.vue";
import OpsHome from "@/pages/DevOps/Home.vue"
import AllDep from "@/pages/DevOps/AllDep.vue";
import MyDep from "@/pages/DevOps/MyDep.vue";
import Issue from "@/pages/DevOps/Issue.vue";
import ManageUsr from "@/pages/DevOps/ManageUsr.vue";
import NewDep from "@/pages/DevOps/NewDep.vue";
import QuickOps from "@/pages/DevOps/QuickOps.vue";
import MacIni from "@/pages/DevOps/MacIni.vue";
import DepIni from "@/pages/DevOps/DepIni.vue";
import DevIni from "@/pages/DevOps/DevIni.vue";
import NewUsr from "@/pages/DevOps/NewUsr.vue";
import DBview from "@/pages/DevOps/DBview.vue";
import DevTools from '@/pages/DevOps/DevTools.vue';
import Note from "@/pages/Note.vue";
import ElkStack from "@/pages/ElkStack.vue";
import Kibana from '@/pages/ElkStack/Kibana.vue';
import IFrame from "@/pages/IFrame.vue";
// import BasicFrame from '@/pages/IFrame/BasicFrame.vue';
import Raft from '@/pages/IFrame/Raft.vue';
import GM239 from '@/pages/IFrame/GM239.vue';
import GM from '@/pages/IFrame/GM.vue';
import Plans from '@/pages/IFrame/Plans.vue'
import ClientHub from '@/pages/IFrame/ClientHub.vue';
import DevDoc from '@/pages/IFrame/DevDoc.vue'
import TaomeeDoc from '@/pages/IFrame/TaomeeDoc.vue'
import Visitor from "@/pages/Visitor.vue";
import BreakDep from "@/pages/Visitor/BreakDep.vue";
import Data from "@/pages/Data.vue";
import e404 from '@/pages/E404.vue';
import UsrInfo from '@/pages/UsrInfo.vue';

Vue.use(Router)

const router = new Router({
    routes: [
        {
            path: "/",
            redirect: "/home"
        },
        {
            path: "/home",
            component: Home
        },
        {
            path: "/devops",
            component: DevOps,
            children: [
                {
                    path: "",
                    redirect: "home"
                },
                {
                    path: "home",
                    component: OpsHome
                },
                {
                    path: "mydep",
                    component: MyDep
                },
                {
                    path: "alldep",
                    component: AllDep
                },
                {
                    path: "manageuser",
                    component: ManageUsr
                },
                {
                    path: "quickops",
                    component: QuickOps
                },
                {
                    path: "newdep",
                    component: NewDep
                },
                {
                    path: "macini",
                    component: MacIni
                },
                {
                    path: "depini",
                    component: DepIni
                },
                {
                    path: "devini",
                    component: DevIni
                },
                {
                    path: "newusr",
                    component: NewUsr
                },
                {
                    path: "dbview",
                    component: DBview
                },
                {
                    path: "devtools",
                    component: DevTools
                },
                {
                    path: "issue",
                    component: Issue
                }
            ]
        },
        {
            path: "/visitor",
            component: Visitor,
            children: [
                {
                    path: "breakdep",
                    component: BreakDep
                },
            ]
        },
        {
            path: "/elkstack",
            component: ElkStack,
            children: [
                {
                    path: "kibana",
                    component: Kibana
                },
            ]
        },
        {
            path: "/note",
            component: Note
        },
        {
            path: "/iframe",
            component: IFrame,
            children: [
                // {
                //     path: "basic",
                //     component: BasicFrame
                // },
                {
                    path: "raft",
                    component: Raft
                },
                {
                    path: "plans",
                    component: Plans
                },
                {
                    path: "gm239",
                    component: GM239
                },
                {
                    path: "gm",
                    component: GM
                },
                {
                    path: "devdoc",
                    component: DevDoc
                },
                {
                    path: "taomeedoc",
                    component: TaomeeDoc
                },
                {
                    path: "clienthub",
                    component: ClientHub
                }
            ]
        },
        {
            path: "/data",
            component: Data
        },
        {
            path: "/usrinfo",
            component: UsrInfo
        },
        {
            path: '*',
            component: e404
        }
    ]
})
export default router