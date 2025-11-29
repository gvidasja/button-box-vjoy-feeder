import { Events, WML } from '@wailsio/runtime'
import { useEffect, useRef, useState } from 'react'

type HanbrakeState = {
  min: number
  max: number
  state: number
}

const BUTTONS = [
  [[1, 2], [9], [13], [17], [21, 22]],
  [[3, 4], [10], [14], [18], [23, 24]],
  [[5, 6], [11], [15], [19], [25, 26]],
  [[7, 8], [12], [16], [20], [27, 28]],
]

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
          <div>Hanbrake</div>
          <div>
            {hanbrake.min} {hanbrake.state} {hanbrake.max}
          </div>
          <div>0 {(hanbrake.state / hanbrake.max).toFixed(2)} 1</div>
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
          <div></div>
          {BUTTONS.map(buttonRow => (
            <div style={{ display: 'flex', gap: 2 }}>
              {buttonRow.map(buttonGroup => (
                <div>
                  {buttonGroup.map(btn => (
                    <Button index={btn} active={btn === button} />
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
  const visualActive = useRef(0)

  useEffect(() => {
    if (active) {
      visualActive.current = visualActive.current + 1
      setTimeout(() => (visualActive.current = visualActive.current - 1), 1000)
    }
  }, [active])

  return (
    <div
      style={{
        display: 'inline-block',
        border: '1px solid black',
        textAlign: 'center',
        width: 20,
        color: visualActive.current ? 'white' : 'black',
        background: visualActive.current ? 'red' : 'white',
      }}
    >
      {index} - {visualActive.current}
    </div>
  )
}

export default App
