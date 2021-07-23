<template>
  <div>
    <div class="columns">
      <div class="column">
        <b-loading :is-full-page="true" :active.sync="isLoading"></b-loading>
        <div v-if="items.length > 0 && !isLoading">
          <b-table :data="items" ref="table" backend-sorting :default-sort="[sortField, sortOrder]" @sort="onSort"
                   :paginated="true">
            <b-table-column style="word-break:break-all;" class="is-one-fifth" field="filename" :label="$t('File')"
                            sortable v-slot="props">
              {{ props.row.filename }}
              <br/><small>{{ props.row.path }}</small>
            </b-table-column>
            <b-table-column field="created_time" :label="$t('Created')" style="white-space: nowrap;" sortable
                            v-slot="props">
              {{ format(parseISO(props.row.created_time), "yyyy-MM-dd hh:mm:ss") }}
            </b-table-column>
            <b-table-column field="size" :label="$t('Size')" style="white-space: nowrap;" sortable v-slot="props">
              {{ prettyBytes(props.row.size) }}
            </b-table-column>
            <b-table-column field="video_width" :label="$t('Width')" sortable v-slot="props">
              <span v-if="props.row.video_width !== 0">{{ props.row.video_width }}</span>
              <span v-else>-</span>
            </b-table-column>
            <b-table-column field="video_height" :label="$t('Height')" sortable v-slot="props">
              <span v-if="props.row.video_height !== 0">{{ props.row.video_height }}</span>
              <span v-else>-</span>
            </b-table-column>
            <b-table-column field="video_bitrate" :label="$t('Bitrate')" style="white-space: nowrap;" sortable
                            v-slot="props">
              <span v-if="props.row.video_bitrate !== 0">{{ prettyBytes(props.row.video_bitrate) }}</span>
              <span v-else>-</span>
            </b-table-column>
            <b-table-column field="duration" :label="$t('Duration')" style="white-space: nowrap;" sortable
                            v-slot="props">
              <span v-if="props.row.duration !== 0">{{ humanizeSeconds(props.row.duration) }}</span>
              <span v-else>-</span>
            </b-table-column>
            <b-table-column field="video_avgfps_val" :label="$t('FPS')" style="white-space: nowrap;" sortable
                            v-slot="props">
              <span v-if="props.row.video_avgfps_val !== 0">{{ props.row.video_avgfps_val }}</span>
              <span v-else>-</span>
            </b-table-column>
            <b-table-column v-slot="props">
              <div class="block">
                <b-button @click="play(props.row)" v-if="props.row.type === 'video'">{{ $t('Play') }}</b-button>
                <b-button v-if="props.row.scene_id === 0" @click="match(props.row)">{{ $t('Match') }}</b-button>
                <b-button v-else @click="showDetails(props.row)">Scene</b-button>
              </div>
            </b-table-column>
            <b-table-column style="white-space: nowrap;" v-slot="props">
              <button class="button is-danger is-outlined" @click='removeFile(props.row)' title='Delete file from disk'>
                <b-icon pack="fas" icon="trash"></b-icon>
              </button>
            </b-table-column>
          </b-table>
        </div>
        <div v-if="items.length === 0 && !isLoading">
          <section class="hero is-large">
            <div class="hero-body">
              <div class="container has-text-centered">
                <h1 class="title">
                  <span class="icon">
                    <i class="far fa-check-circle is-superlarge"></i>
                  </span>
                </h1>
                <h2 class="subtitle">
                  {{ $t('No files matching your selection') }}
                </h2>
              </div>
            </div>
          </section>
        </div>
      </div>
    </div>
    <Details v-if="showOverlay"/>
    <EditScene v-if="showEdit" />
  </div>
</template>

<script>
import prettyBytes from 'pretty-bytes'
import { format, parseISO } from 'date-fns'
import ky from 'ky'
import Details from '../scenes/Details'
import EditScene from '../scenes/EditScene'

export default {
  name: 'List',
  components: { Details, EditScene },
  data () {
    return {
      files: [],
      prettyBytes,
      format,
      parseISO,
      sortField: 'created_time',
      sortOrder: 'desc'
    }
  },
  computed: {
    isLoading () {
      return this.$store.state.files.isLoading
    },
    items () {
      return this.$store.state.files.items
    },
    showOverlay() {
      return this.$store.state.overlay.details.show
    },
    showEdit() {
      return this.$store.state.overlay.edit.show
    }
  },
  mounted () {
    this.$store.state.files.filters.sort = `${this.sortField}_${this.sortOrder}`
    this.$store.dispatch('files/load')
  },
  methods: {
    onSort (field, order) {
      this.sortField = field
      this.sortOrder = order
      this.$store.state.files.filters.sort = `${field}_${order}`
      this.$store.dispatch('files/load')
    },
    play (file) {
      ky.get(`/deovr/local/${encodeURIComponent(`${file.path}\\${file.filename}`)}`);
    },
    async showDetails(file) {
      const scene = await fetch("http://localhost:9999/api/scene/list", {
        "headers": {
          "content-type": "application/json"
        },
        "body": JSON.stringify({id:file.scene_id}),
        "method": "POST"
      }).then(res => res.json()).then(res => res.scenes.pop());
      this.$store.commit("overlay/showDetails", {scene: scene});
    },
    match (file) {
      this.$store.commit('overlay/showMatch', { file: file })
    },
    humanizeSeconds (seconds) {
      return new Date(seconds * 1000).toISOString().substr(11, 8)
    },
    removeFile (file) {
      this.$buefy.dialog.confirm({
        title: 'Remove file',
        message: `You're about to remove file <strong>${file.filename}</strong> from <strong>disk</strong>.`,
        type: 'is-danger',
        hasIcon: true,
        onConfirm: () => {
          ky.delete(`/api/files/file/${file.id}`).json().then(data => {
            this.$store.dispatch('files/load')
          })
        }
      })
    }
  }
}
</script>

<style scoped>
small {
  opacity: 0.6;
}

.is-superlarge {
  height: 96px;
  max-height: 96px;
  max-width: 96px;
  min-height: 96px;
  min-width: 96px;
  width: 96px;
}
</style>
