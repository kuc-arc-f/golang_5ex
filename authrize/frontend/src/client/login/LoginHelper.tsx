import LibConfig from "../lib/LibConfig"
import {SendMessage} from "../../../wailsjs/go/main/App";

const LoginHelper = {

  user_login: (email , password) => {
    try{    
      const inData = {
        email: email,
        password: password
      };        
      const target = {
        action: "user_get",
        data: JSON.stringify(inData)
      }        
      const sendJson = JSON.stringify(target)        
      console.log(sendJson)  
      SendMessage(sendJson).then((result) => {
        console.log("result=", result);
        const j1 = JSON.parse(result)
        console.log(j1);
        if(j1.Ret === 200){
          alert("Succes , Login");
          localStorage.setItem(LibConfig.STORAGE_KEY_USER_ID, 1);
        }
        if(j1.Ret === 400){
          alert("NG , user data");
        }
      }).catch((err) => { console.error(err);});        
    } catch (error) {
      console.error('Error fetching:', error);
    }
  },

}
export default LoginHelper;
