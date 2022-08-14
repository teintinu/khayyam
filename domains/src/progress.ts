import { asap } from 'pjobs'

export interface Progress {
    readonly actual: Status
    on(callback: (status: Status) => void): () => void
    update(status: Status): Promise<void>
}

export type State = 'manual' | 'created' | 'idle' | 'waiting'| 'working' | 'failed' | 'finished' | 'killed'

export interface Status {
    readonly state: State
    readonly message?: string
    readonly detail?: string
    readonly percent?: number
}

export function createProgress (manual: boolean): Progress {
  const callbacks = new Set<(status: Status) => void>()
  let actual: Status = { state: manual ? 'manual' : 'created' }
  return {
    get actual () {
      return actual
    },
    on (callback) {
      callbacks.add(callback)
      return () => {
        callbacks.delete(callback)
      }
    },
    async update (status) {
      const changed = actual.state !== status.state ||
        actual.message !== status.message ||
        actual.detail !== status.detail ||
        actual.percent !== status.percent
      if (changed) {
        actual = status
        const promises: Array<Promise<void>> = []
        callbacks.forEach((callback) => {
          promises.push(new Promise((resolve) => {
            asap(() => {
              try {
                callback(status)
              } catch (e) {
              // ignore
              }
              resolve()
            })
          }))
        })
        await Promise.all(promises)
      }
    }
  }
}
