import { Events, WML } from '@wailsio/runtime'
import { useEffect, useState } from 'react'

type HanbrakeState = {
  min: number
  max: number
  state: number
}

function App() {
  const [hanbrake, setHanbrake] = useState<HanbrakeState>({ min: 0, max: 1, state: 0 })
  const [button, setButton] = useState<number>(0)

  useEffect(() => {
    Events.On('handbrake', (v: Events.WailsEvent) => {
      setHanbrake(v.data[0])
    })

    Events.On('button', (v: Events.WailsEvent) => {
      setButton(v.data[0])
    })

    WML.Reload()
  }, [])

  return (
    <div>
      <div>
        <div>
          <span>Hanbrake</span>
          <span>
            {hanbrake.min} {hanbrake.state} {hanbrake.max}
          </span>
        </div>
        <div>
          <progress max={hanbrake.max} value={hanbrake.state - hanbrake.min}></progress>
        </div>
      </div>
      <div>
        <div>
          <span>Button:</span>
          <span>{button}</span>
        </div>
        <div>
          {Array.from({ length: 32 }, (_, i) => (
            <div
              key={i}
              style={{
                display: 'inline-block',
                border: '1px solid black',
                textAlign: 'center',
                width: 20,
                background: button === i ? 'red' : 'white',
              }}
            >
              {i + 1}
            </div>
          ))}
        </div>
      </div>
    </div>
  )
}

export default App
