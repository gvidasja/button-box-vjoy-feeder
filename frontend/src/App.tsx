import { Events, WML } from '@wailsio/runtime'
import { useEffect, useState } from 'react'

type HanbrakeState = {
  min: number
  max: number
  state: number
}

const BUTTONS = [
  [[1, 2], [9, 13, 17], [21, 22]],
  [[3, 4], [10, 14, 18], [23, 24]],
  [[5, 6], [11, 15, 19], [25, 26]],
  [[7, 8], [12, 16, 20], [27, 28]],
]

function App() {
  const [handbrake, setHanbrake] = useState<HanbrakeState>({ min: 0, max: 1, state: 0 })
  const [buttons, setButtons] = useState<Record<number, boolean>>({})
  const [lastButtonEvent, setLastButtonEvent] = useState<{ button: number; state: boolean } | null>(null)

  useEffect(() => {
    Events.On('handbrake', (v: Events.WailsEvent) => {
      setHanbrake(v.data)
    })

    Events.On('button', (v: Events.WailsEvent) => {
      setButtons({ ...buttons, [v.data.button]: v.data.state })
      setLastButtonEvent({ button: v.data.button, state: v.data.state })
    })

    WML.Reload()
  }, [])

  return (
    <div>
      <div>
        <div>
          <div>Hanbrake</div>
          <div>
            {handbrake.min} {handbrake.state} {handbrake.max}
          </div>
          <div>0 {(handbrake.state / handbrake.max).toFixed(2)} 1</div>
        </div>
        <div>
          <progress max={handbrake.max} value={handbrake.state - handbrake.min}></progress>
        </div>
      </div>
      <div>
        <div>
          <div>Last event: {lastButtonEvent ? `Button ${lastButtonEvent.button} - ${lastButtonEvent.state}` : 'None'}</div>
          {BUTTONS.map(buttonRow => (
            <div style={{ display: 'flex', gap: 2 }}>
              {buttonRow.map(buttonGroup => (
                <div>
                  {buttonGroup.map(btn => (
                    <Button index={btn} active={buttons[btn]} />
                  ))}
                </div>
              ))}
            </div>
          ))}
        </div>
      </div>
      <div style={{ position: 'fixed', bottom: 0, right: 0 }}>
        Built at: {import.meta.env.BUILD_TIME as string}
      </div>
    </div>
  )
}

function Button({ active, index }: { active: boolean; index: number }) {
  return (
    <div
      style={{
        display: 'inline-block',
        border: '1px solid black',
        textAlign: 'center',
        width: 20,
        color: active ? 'white' : 'black',
        background: active ? 'red' : 'white',
      }}
    >
      {index}
    </div>
  )
}

export default App
