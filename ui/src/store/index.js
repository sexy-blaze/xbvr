import Vue from "vue";
import Vuex from "vuex";

import actorList from "./actorList";
import sceneList from "./sceneList";
import messages from "./messages";
import overlay from "./overlay";
import files from "./files";
import optionsStorage from "./optionsStorage";
import optionsDLNA from "./optionsDLNA";
import optionsSites from "./optionsSites";
import optionsPreviews from "./optionsPreviews";


Vue.use(Vuex);

export default new Vuex.Store({
  modules: {
    actorList,
    sceneList,
    messages,
    overlay,
    files,
    optionsStorage,
    optionsDLNA,
    optionsSites,
    optionsPreviews,
  }
})
