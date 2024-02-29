import './MyInput.scss'

interface MyInputProps {
  title: string
}


function MyInput(props: MyInputProps) {
  return(
    <div className="my-input">
      <span>{props.title}</span>
      <input placeholder={props.title} />
    </div>
  )
}

function Auth() {
  return (
    <div className='page-panel'>
      <h1>Slyfox-Korsak</h1>
      <div className='card'>
        <form>
          <h1>Log In</h1>
          <MyInput title="Username"/>
          <MyInput title="Password"/>
          <button>Log In</button>
        </form>
      </div>
    </div>
  )
}

export default Auth
