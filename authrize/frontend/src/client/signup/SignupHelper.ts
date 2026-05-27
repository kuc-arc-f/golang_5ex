import {SendMessage} from "../../../wailsjs/go/main/App";

const SignupHelper = {

  user_create: (email , password, username) => {
    try{    
      const inData = {
        email: email,
        password: password,
        username: username        
      };        
      const target = {
        action: "user_create",
        data: JSON.stringify(inData)
      }        
      const sendJson = JSON.stringify(target)        
      console.log(sendJson)
      SendMessage(sendJson).then((result) => {
        console.log("result=", result);
        const j1 = JSON.parse(result)
        console.log(j1);
        if(j1.Ret === 200){
          //fetchTodos();
          alert("Succes , send data");
        }
      }).catch((err) => { console.error(err);});         
    } catch (error) {
      console.error('Error fetching:', error);
    }
  },

}
export default SignupHelper;
