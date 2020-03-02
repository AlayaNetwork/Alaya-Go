package network.platon.contracts.wasm;

import com.platon.rlp.datatypes.Uint32;
import java.util.Arrays;
import org.web3j.abi.WasmFunctionEncoder;
import org.web3j.abi.datatypes.WasmFunction;
import org.web3j.crypto.Credentials;
import org.web3j.protocol.Web3j;
import org.web3j.protocol.core.RemoteCall;
import org.web3j.tx.TransactionManager;
import org.web3j.tx.WasmContract;
import org.web3j.tx.gas.GasProvider;

/**
 * <p>Auto generated code.
 * <p><strong>Do not modify!</strong>
 * <p>Please use the <a href="https://docs.web3j.io/command_line.html">web3j command line tools</a>,
 * or the org.web3j.codegen.SolidityFunctionWrapperGenerator in the 
 * <a href="https://github.com/web3j/web3j/tree/master/codegen">codegen module</a> to update.
 *
 * <p>Generated with web3j version 0.9.1.0-SNAPSHOT.
 */
public class Sha3Function extends WasmContract {
    private static String BINARY_0 = "0x0061736d01000000014c0d60000060017f0060027f7f0060047f7f7f7f006000017f60037f7f7f017f60017f017f60027f7f017f60037f7f7f0060047f7f7f7f017f60037f7e7e017f60027e7e017f60047f7e7e7f000282010603656e760c706c61746f6e5f6465627567000203656e760c706c61746f6e5f70616e6963000003656e760b706c61746f6e5f73686133000303656e7617706c61746f6e5f6765745f696e7075745f6c656e677468000403656e7610706c61746f6e5f6765745f696e707574000103656e760d706c61746f6e5f72657475726e0002036867000005050606060501070100060105010108050607010005020706070601060107010506060606030208060106050902070905060303050602020702070207020802080502020a0b06080c020306020301020707050302020709020607030105020001010600000405017001040405030100020615037f0141d08c040b7f0041d08c040b7f0041cf0c0b075406066d656d6f72790200115f5f7761736d5f63616c6c5f63746f727300060b5f5f686561705f6261736503010a5f5f646174615f656e6403020f5f5f66756e63735f6f6e5f65786974001c06696e766f6b65006b0909010041010b032327680ae75767080010111067106c0b02000bc60a010b7f2002410f6a210341002104410020026b21052002410e6a2106410120026b21072002410d6a2108410220026b210902400340200020046a210b200120046a210a20022004460d01200a410371450d01200b200a2d00003a00002003417f6a2103200541016a21052006417f6a2106200741016a21072008417f6a2108200941016a2109200441016a21040c000b0b200220046b210c02400240024002400240200b410371220d450d00200c4120490d03200d4101460d01200d4102460d02200d4103470d03200b200120046a28020022063a0000200041016a210c200220046b417f6a21092004210b0240034020094113490d01200c200b6a220a2001200b6a220741046a2802002208411874200641087672360200200a41046a200741086a2802002206411874200841087672360200200a41086a2007410c6a2802002208411874200641087672360200200a410c6a200741106a2802002206411874200841087672360200200b41106a210b200941706a21090c000b0b2002417f6a2005416d2005416d4b1b20036a4170716b20046b210c2001200b6a41016a210a2000200b6a41016a210b0c030b200c210a02400340200a4110490d01200020046a220b200120046a2207290200370200200b41086a200741086a290200370200200441106a2104200a41706a210a0c000b0b02400240200c4108710d00200120046a210a200020046a21040c010b200020046a220b200120046a2204290200370200200441086a210a200b41086a21040b0240200c410471450d002004200a280200360200200a41046a210a200441046a21040b0240200c410271450d002004200a2f00003b0000200441026a2104200a41026a210a0b200c410171450d032004200a2d00003a000020000f0b200b200120046a220a28020022063a0000200b41016a200a41016a2f00003b0000200041036a210c200220046b417d6a21052004210b0240034020054111490d01200c200b6a220a2001200b6a220741046a2802002203410874200641187672360200200a41046a200741086a2802002206410874200341187672360200200a41086a2007410c6a2802002203410874200641187672360200200a410c6a200741106a2802002206410874200341187672360200200b41106a210b200541706a21050c000b0b2002417d6a2009416f2009416f4b1b20086a4170716b20046b210c2001200b6a41036a210a2000200b6a41036a210b0c010b200b200120046a220a28020022083a0000200b41016a200a41016a2d00003a0000200041026a210c200220046b417e6a21052004210b0240034020054112490d01200c200b6a220a2001200b6a220941046a2802002203411074200841107672360200200a41046a200941086a2802002208411074200341107672360200200a41086a2009410c6a2802002203411074200841107672360200200a410c6a200941106a2802002208411074200341107672360200200b41106a210b200541706a21050c000b0b2002417e6a2007416e2007416e4b1b20066a4170716b20046b210c2001200b6a41026a210a2000200b6a41026a210b0b0240200c411071450d00200b200a2d00003a0000200b200a280001360001200b200a290005370005200b200a2f000d3b000d200b200a2d000f3a000f200b41106a210b200a41106a210a0b0240200c410871450d00200b200a290000370000200b41086a210b200a41086a210a0b0240200c410471450d00200b200a280000360000200b41046a210b200a41046a210a0b0240200c410271450d00200b200a2f00003b0000200b41026a210b200a41026a210a0b200c410171450d00200b200a2d00003a00000b20000bfb0202027f017e02402002450d00200020013a0000200020026a2203417f6a20013a000020024103490d00200020013a0002200020013a00012003417d6a20013a00002003417e6a20013a000020024107490d00200020013a00032003417c6a20013a000020024109490d002000410020006b41037122046a2203200141ff017141818284086c22013602002003200220046b417c7122046a2202417c6a200136020020044109490d002003200136020820032001360204200241786a2001360200200241746a200136020020044119490d002003200136021820032001360214200320013602102003200136020c200241706a20013602002002416c6a2001360200200241686a2001360200200241646a20013602002001ad220542208620058421052004200341047141187222016b2102200320016a2101034020024120490d0120012005370300200141186a2005370300200141106a2005370300200141086a2005370300200141206a2101200241606a21020c000b0b20000b7a01027f200021010240024003402001410371450d0120012d0000450d02200141016a21010c000b0b2001417c6a21010340200141046a22012802002202417f73200241fffdfb776a7141808182847871450d000b0340200241ff0171450d01200141016a2d00002102200141016a21010c000b0b200120006b0b3a01017f23808080800041106b220141d08c84800036020c2000200128020c41076a41787122013602042000200136020020003f0036020c20000b120041808880800020004108108d808080000bc70101067f23808080800041106b22032480808080002003200136020c024002402001450d002000200028020c200241036a410020026b220471220520016a220641107622076a220836020c200020022000280204220120066a6a417f6a20047122023602040240200841107420024b0d002000410c6a200841016a360200200741016a21070b0240200740000d0041c88a808000108e808080000b20012003410c6a41041088808080001a200120056a21000c010b410021000b200341106a24808080800020000b180020002000108a808080001080808080001081808080000b2e000240418088808000200120006c22004108108d808080002201450d002001410020001089808080001a0b20010b02000b0f00418088808000108b808080001a0b3a01027f2000410120001b2101024003402001108c8080800022020d014100280290888080002200450d012000118080808000000c000b0b20020b0a0020001090808080000bce0301067f024020002001460d000240024002400240200120006b20026b410020024101746b4d0d0020012000734103712103200020014f0d012003450d02200021030c030b2000200120021088808080000f0b024020030d002001417f6a210402400340200020026a2203410371450d012002450d052003417f6a200420026a2d00003a00002002417f6a21020c000b0b2000417c6a21032001417c6a2104034020024104490d01200320026a200420026a2802003602002002417c6a21020c000b0b2001417f6a210103402002450d03200020026a417f6a200120026a2d00003a00002002417f6a21020c000b0b200241046a21052002417f7321064100210402400340200120046a2107200020046a2208410371450d0120022004460d03200820072d00003a00002005417f6a2105200641016a2106200441016a21040c000b0b200220046b2101410021030240034020014104490d01200820036a200720036a280200360200200341046a21032001417c6a21010c000b0b200720036a2101200820036a210320022006417c2006417c4b1b20056a417c716b20046b21020b03402002450d01200320012d00003a00002002417f6a2102200341016a2103200141016a21010c000b0b20000b0900200041013602000b0900200041003602000b7701027f0240200241704f0d00024002402002410a4b0d00200020024101743a0000200041016a21030c010b200241106a417071220410928080800021032000200236020420002004410172360200200020033602080b2003200120021098808080001a200320026a41003a00000f0b108780808000000b1a0002402002450d0020002001200210888080800021000b20000b1d00024020002d0000410171450d0020002802081093808080000b20000b270020004200370200200041086a4100360200200020012001108a8080800010978080800020000b0900108780808000000bb60101037f4194888080001095808080004100280298888080002100024003402000450d01024003404100410028029c888080002202417f6a220136029c8880800020024101480d01200020014102746a22004184016a2802002102200041046a2802002100419488808000109680808000200220001181808080000041948880800010958080800041002802988880800021000c000b0b4100412036029c88808000410020002802002200360298888080000c000b0b0bcd0101027f419488808000109580808000024041002802988880800022030d0041a0888080002103410041a088808000360298888080000b02400240410028029c8880800022044120470d004184024101108f808080002203450d0141002104200341002802988880800036020041002003360298888080004100410036029c888080000b4100200441016a36029c88808000200320044102746a22034184016a2001360200200341046a200036020041948880800010968080800041000f0b419488808000109680808000417f0b6001017f23808080800041206b2202248080808000200241186a420037030020024200370310200242003703082000200241086a2001109f8080800010a08080800010a1808080001a200241086a10a2808080001a200241206a2480808080000b4101017f23808080800041106b2202248080808000200020022001109a80808000220110d880808000210020011099808080001a200241106a24808080800020000b23000240200028020c200041106a280200460d0041e98b808000108e808080000b20000b4e01017f20004200370200200041003602080240200128020420012802006b2202450d002000200210d180808000200041086a2001280200200141046a280200200041046a10d2808080000b20000b19002000410c6a10d3808080001a200010a4808080001a20000b0f0041a48a80800010a4808080001a0b2201017f024020002802002201450d002000200136020420011093808080000b20000b4701027f23808080800041206b22012480808080002000200141086a410010a680808000220210a08080800010a1808080001a200210a2808080001a200141206a2480808080000b24002000420037020820004200370200200041106a42003702002000200110c0808080000b0f0041b08a80800010a4808080001a0b95010020004200370210200042ffffffff0f3702082000200129020037020002402002410871450d00200010a98080800020012802044f0d00024020024104710d00200042003702000c010b41e28a808000108e808080000b024002402002411071450d00200010a98080800020012802044d0d0020024104710d01200042003702000b20000f0b41f08a808000108e8080800020000b3400024002402000280204450d0020002802002c0000417f4c0d0141010f0b41000f0b200010aa80808000200010ab808080006a0b280002402000280204450d0020002802002c0000417f4c0d0041000f0b200010b08080800041016a0b980401047f0240024002402000280204450d00200010b1808080004101210120002802002c00002202417f4c0d010c020b41000f0b0240200241ff0171220141b7014b0d00200141807f6a0f0b02400240200241ff0171220241bf014b0d000240200041046a22032802002202200141c97e6a22044b0d0041ff8a808000108e80808000200328020021020b024020024102490d0020002802002d00010d0041ff8a808000108e808080000b024020044105490d0041f08a808000108e808080000b024020002802002d00010d0041ff8a808000108e808080000b41002101410021020240034020042002460d012001410874200028020020026a41016a2d0000722101200241016a21020c000b0b200141384f0d0141ff8a808000108e8080800020010f0b0240200241f7014b0d00200141c07e6a0f0b0240200041046a22032802002202200141897e6a22044b0d0041ff8a808000108e80808000200328020021020b024020024102490d0020002802002d00010d0041ff8a808000108e808080000b024020044105490d0041f08a808000108e808080000b024020002802002d00010d0041ff8a808000108e808080000b41002101410021020240034020042002460d012001410874200028020020026a41016a2d0000722101200241016a21020c000b0b200141384f0d0041ff8a808000108e8080800020010f0b200141ff7d490d0041f08a808000108e8080800020010f0b20010b5102017f017e23808080800041306b220124808080800020012000290200220237031020012002370308200141186a200141086a411410a88080800010a9808080002100200141306a24808080800020000b6a01037f02400240024020012802002204450d0041002105200320026a200128020422064b0d0120062002490d014100210120062003490d02200620026b20032003417f461b2101200420026a21050c020b410021050b410021010b20002001360204200020053602000b3901017f0240200110ab80808000220220012802044d0d0041808c808000108e808080000b20002001200110aa80808000200210ad808080000bd003020a7f017e23808080800041c0006b220324808080800002402001280208220420024d0d00200341386a200110ae80808000200320032903383703182001200341186a10ac8080800036020c200341306a200110ae80808000410021044100210541002106024020032802302207450d00410021054100210620032802342208200128020c2209490d00200820092009417f461b2105200721060b20012006360210200141146a2005360200200141086a41003602000b200141106a2106200141146a21092001410c6a2107200141086a210802400340200420024f0d012009280200450d01200341306a200110ae8080800041002104024002402003280230220a450d00410021052003280234220b2007280200220c490d01200a200c6a2105200b200c6b21040c010b410021050b20092004360200200620053602002003200436022c2003200536022820032003290328370310200341306a20064100200341106a10ac8080800010ad8080800020062003290330220d37020020072007280200200d422088a76a3602002008200828020041016a22043602000c000b0b20032006290200220d3703202003200d3703082000200341086a411410a8808080001a200341c0006a2480808080000b4701017f4100210102402000280204450d00024020002802002d0000220041bf014b0d00200041b801490d01200041c97e6a0f0b200041f801490d00200041897e6a21010b20010b6601017f024020002802040d0041ff8a808000108e808080000b0240200028020022012d0000418101470d000240200041046a28020041014b0d0041ff8a808000108e80808000200028020021010b20012c00014100480d0041ff8a808000108e808080000b0bc80101047f024002402000280204450d00200010b180808000200028020022012c000022024100480d0120024100470f0b41000f0b410121030240200241807f460d000240200241ff0171220441b7014b0d000240200041046a28020041014b0d0041ff8a808000108e80808000200028020021010b20012d00014100470f0b41002103200441bf014b0d000240200041046a280200200241ff017141ca7e6a22024b0d0041ff8a808000108e80808000200028020021010b200120026a2d000041004721030b20030b2d01017f2000200028020420012802002203200320012802046a10b4808080001a2000200210b58080800020000b970201057f23808080800041206b22042480808080000240200320026b22054101480d00024020052000280208200028020422066b4c0d00200441086a2000200520066a20002802006b10b680808000200120002802006b200041086a10b78080800021060240034020032002460d01200641086a220528020020022d00003a00002005200528020041016a360200200241016a21020c000b0b20002006200110b8808080002101200610b9808080001a0c010b024002402005200620016b22074c0d00200041086a200220076a22082003200041046a10ba80808000200741014e0d010c020b200321080b200020012006200120056a10bb8080800020022008200110bc808080001a0b200441206a24808080800020010bd00201087f02402001450d002000410c6a2102200041106a2103200041046a21040340200328020022052002280200460d010240200541786a28020020014f0d0041878b808000108e80808000200328020021050b200541786a2206200628020020016b220136020020010d0120032006360200200428020020002802006b2005417c6a28020022016b220510bd8080800021062000200428020020002802006b22074101200641016a20054138491b22086a10be80808000200120002802006a220920086a2009200720016b1094808080001a02400240200541374b0d00200028020020016a200541406a3a00000c010b0240200641f7016a220741ff014b0d00200028020020016a20073a00002000280200200620016a6a210103402005450d02200120053a0000200541087621052001417f6a21010c000b0b419b8b808000108e808080000b410121010c000b0b0b4c01017f02402001417f4c0d0041ffffffff0721020240200028020820002802006b220041feffffff034b0d0020012000410174220020002001491b21020b20020f0b2000109b80808000000b5401017f410021042000410036020c200041106a200336020002402001450d00200110928080800021040b200020043602002000200420026a22023602082000410c6a200420016a3602002000200236020420000b8c0101027f20012802042103200041086a220420002802002002200141046a10d580808000200420022000280204200141086a10db80808000200028020021022000200128020436020020012002360204200028020421022000200128020836020420012002360208200028020821022000200128020c3602082001200236020c2001200128020436020020030b2301017f200010d680808000024020002802002201450d0020011093808080000b20000b2e000240200220016b22024101480d002003280200200120021088808080001a2003200328020020026a3602000b0b5c01037f200041046a21042000280204220521062001200520036b6a2203210002400340200020024f0d01200620002d00003a00002004200428020041016a2206360200200041016a21000c000b0b20012003200510da808080001a0b21000240200120006b2201450d002002200020011094808080001a0b200220016a0b2501017f41002101024003402000450d0120004108762100200141016a21010c000b0b20010b4001027f02402000280204200028020022026b220320014f0d002000200120036b10bf808080000f0b0240200320014d0d00200041046a200220016a3602000b0b920101027f23808080800041206b2202248080808000024002402000280208200028020422036b20014f0d00200241086a2000200320016a20002802006b10b680808000200041046a28020020002802006b200041086a10b7808080002203200110dc808080002000200310d480808000200310b9808080001a0c010b2000200110dd808080000b200241206a2480808080000b7501017f23808080800041106b2202248080808000024002402001450d00200220013602002002200028020420002802006b3602042000410c6a200210c1808080000c010b20024100360208200242003703002000200210c2808080001a200210a4808080001a0b200241106a24808080800020000b3d01017f02402000280204220220002802084f0d0020022001290200370200200041046a2200200028020041086a3602000f0b2000200110c3808080000b5101027f23808080800041106b22022480808080002002200128020022033602082002200128020420036b36020c200220022903083703002000200210c4808080002101200241106a24808080800020010b840101027f23808080800041206b2202248080808000200241086a2000200028020420002802006b41037541016a10de80808000200028020420002802006b410375200041086a10df80808000220328020820012902003702002003200328020841086a3602082000200310e080808000200310e1808080001a200241206a2480808080000b800102027f017e23808080800041206b2202248080808000024002402001280204220341374b0d002002200341406a3a001f20002002411f6a10c5808080000c010b2000200341f70110c6808080000b200220012902002204370310200220043703082000200241086a410110b3808080002100200241206a24808080800020000b3d01017f02402000280204220220002802084f0d00200220012d00003a0000200041046a2200200028020041016a3602000f0b2000200110c7808080000b6401027f23808080800041106b22032480808080000240200110bd80808000220420026a2202418002480d0041d18b808000108e808080000b200320023a000f20002003410f6a10c58080800020002001200410c880808000200341106a2480808080000b7e01027f23808080800041206b2202248080808000200241086a2000200028020441016a20002802006b10b680808000200028020420002802006b200041086a10b780808000220328020820012d00003a00002003200328020841016a3602082000200310d480808000200310b9808080001a200241206a2480808080000b44002000200028020420026a20002802006b10be808080002000280204417f6a2100024003402001450d01200020013a00002000417f6a2100200141087621010c000b0b0bfc0101037f23808080800041206b22032480808080002001280200210420012802042105024002402002450d004100210102400340200420016a2102200120054f0d0120022d00000d01200141016a21010c000b0b200520016b21050c010b200421020b0240024002400240024020054101470d0020022c00004100480d012000200210ca808080000c040b200541374b0d010b20032005418001733a001f20002003411f6a10c5808080000c010b2000200541b70110c6808080000b2003200536021420032002360210200320032903103703082000200341086a410010b3808080001a0b2000410110b580808000200341206a24808080800020000b3d01017f0240200028020422022000280208460d00200220012d00003a0000200041046a2200200028020041016a3602000f0b2000200110cb808080000b7e01027f23808080800041206b2202248080808000200241086a2000200028020441016a20002802006b10b680808000200028020420002802006b200041086a10b780808000220328020820012d00003a00002003200328020841016a3602082000200310d480808000200310b9808080001a200241206a2480808080000bfa0101047f23808080800041106b220324808080800002400240024020012002844200510d00200142ff005620024200522002501b0d01200320013c000f20002003410f6a10c5808080000c020b200041b78b80800010ca808080000c010b024002402001200210cd80808000220441374b0d0020032004418001733a000e20002003410e6a10c5808080000c010b0240200410ce80808000220541b7016a2206418002490d0041b88b808000108e808080000b200320063a000d20002003410d6a10c58080800020002004200510cf808080000b200020012002200410d0808080000b2000410110b580808000200341106a24808080800020000b3501017f41002102024003402000200184500d0120004208882001423886842100200241016a2102200142088821010c000b0b20020b2501017f41002101024003402000450d0120004108762100200141016a21010c000b0b20010b44002000200028020420026a20002802006b10be808080002000280204417f6a2100024003402001450d01200020013a00002000417f6a2100200141087621010c000b0b0b54002000200028020420036a20002802006b10be808080002000280204417f6a2100024003402001200284500d01200020013c0000200142088820024238868421012000417f6a2100200242088821020c000b0b0b3801017f02402001417f4c0d00200020011092808080002202360200200020023602042000200220016a3602080f0b2000109b80808000000b2e000240200220016b22024101480d002003280200200120021088808080001a2003200328020020026a3602000b0b2201017f024020002802002201450d002000200136020420011093808080000b20000b7001017f200041086a20002802002000280204200141046a10d580808000200028020021022000200128020436020020012002360204200028020421022000200128020836020420012002360208200028020821022000200128020c3602082001200236020c200120012802043602000b2f01017f20032003280200200220016b22026b2204360200024020024101480d002004200120021088808080001a0b0b0f002000200028020410d7808080000b2d01017f20002802082102200041086a21000240034020012002460d0120002002417f6a22023602000c000b0b0b4501017f23808080800041106b22022480808080002002200241086a200110d98080800029020037030020002002410010c9808080002100200241106a24808080800020000b360020002001280208200141016a20012d00004101711b3602002000200128020420012d0000220141017620014101711b36020420000b23000240200120006b2201450d00200220016b2202200020011094808080001a0b20020b2e000240200220016b22024101480d002003280200200120021088808080001a2003200328020020026a3602000b0b3401017f20002802082102200041086a21000340200241003a00002000200028020041016a22023602002001417f6a22010d000b0b3401017f20002802042102200041046a21000340200241003a00002000200028020041016a22023602002001417f6a22010d000b0b5301017f024020014180808080024f0d0041ffffffff0121020240200028020820002802006b220041037541feffffff004b0d0020012000410275220020002001491b21020b20020f0b2000109b80808000000b5c01017f410021042000410036020c200041106a200336020002402001450d002003200110e28080800021040b200020043602002000200420024103746a22033602082000410c6a200420014103746a3602002000200336020420000b7001017f200041086a20002802002000280204200141046a10e380808000200028020021022000200128020436020020012002360204200028020421022000200128020836020420012002360208200028020821022000200128020c3602082001200236020c200120012802043602000b2301017f200010e480808000024020002802002201450d0020011093808080000b20000b0e0020002001410010e5808080000b2f01017f20032003280200200220016b22026b2204360200024020024101480d002004200120021088808080001a0b0b0f002000200028020410e6808080000b2300024020014180808080024f0d0020014103741092808080000f0b108780808000000b2d01017f20002802082102200041086a21000240034020012002460d012000200241786a22023602000c000b0b0b4a0041a48a80800041e18a808000109e808080004181808080004100418088808000109d808080001a41b08a80800010a5808080004182808080004100418088808000109d808080001a0b0f0041bc8a8080001099808080001a0b180020002000108a808080001080808080001081808080000b7b01027f23808080800041306b22012480808080002001412c6a41002d008d8c8080003a0000200141186a4200370300200141106a42003703002001420037030820014200370300200141002800898c808000360228200141286a41052001412010828080800020012802002102200141306a24808080800020020bd80502077f017e23808080800041e0006b22002480808080001086808080001083808080002201108c8080800022021084808080002000200136024c2000200236024820002000290348370308200041106a200041286a200041086a411c10a880808000410010af80808000200041106a10b18080800002400240200041106a10b280808000450d002000280214450d0020002802102d000041c001490d010b41af8c80800010e9808080000b0240200041106a10ab808080002203200028021422044d0d0041af8c80800010e980808000200028021421040b200028021021050240024002400240024002402004450d004100210620052c00002201417f4a0d03200141ff0171220641bf014b0d0141002102200141ff017141b801490d02200641c97e6a21020c020b4101210620050d02410021010c030b41002102200141ff017141f801490d00200641897e6a21020b200241016a21060b41002101200620036a20044b0d0020042003490d004100210220042006490d01200520066a2101200420066b20032003417f461b22024109490d0141af8c80800010e9808080000c010b410021020b4200210702402002450d00034020074208862001310000842107200141016a21012002417f6a22020d000b0b024002400240200742cfc1bbf0fb98b4ed4e510d0020074283bb9aa3ed8ffada69510d0220074200520d01418e8c80800010e9808080000c020b200041c7006a10ea808080002101200041d8006a22024200370300200041d0006a420037030020004200370348200041c8006a2001ad420010cc808080001a024020002802542002280200460d0041b88c80800010e9808080000b20002802482201200028024c20016b1085808080000240200041d4006a2802002201450d00200041d8006a200136020020011093808080000b20002802482201450d012000200136024c20011093808080000c010b419c8c80800010e9808080000b200041e0006a2480808080000b5501017f410042003702bc8a808000410041003602c48a8080004174210002404174450d000340200041c88a8080006a4100360200200041046a22000d000b0b4183808080004100418088808000109d808080001a0b0bde0402004180080bc802000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000041c80a0b87026661696c656420746f20616c6c6f6361746520706167657300006f7665722073697a6520726c7000756e6465722073697a6520726c700062616420726c70006974656d436f756e7420746f6f206c61726765006974656d436f756e7420746f6f206c6172676520666f7220524c5000804e756d62657220746f6f206c6172676520666f7220524c5000436f756e7420746f6f206c6172676520666f7220524c50006c697374537461636b206973206e6f7420656d70747900626164206361737400010710203076616c6964206d6574686f640a006e6f206d6574686f6420746f2063616c6c0a006261642063617374006c697374537461636b206973206e6f7420656d70747900";

    public static String BINARY = BINARY_0;

    public static final String FUNC_SHA3RESULT = "Sha3Result";

    protected Sha3Function(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, credentials, contractGasProvider);
    }

    protected Sha3Function(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        super(BINARY, contractAddress, web3j, transactionManager, contractGasProvider);
    }

    public static RemoteCall<Sha3Function> deploy(Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(Sha3Function.class, web3j, credentials, contractGasProvider, encodedConstructor);
    }

    public static RemoteCall<Sha3Function> deploy(Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        String encodedConstructor = WasmFunctionEncoder.encodeConstructor(BINARY, Arrays.asList());
        return deployRemoteCall(Sha3Function.class, web3j, transactionManager, contractGasProvider, encodedConstructor);
    }

    public RemoteCall<Uint32> Sha3Result() {
        final WasmFunction function = new WasmFunction(FUNC_SHA3RESULT, Arrays.asList(), Uint32.class);
        return executeRemoteCall(function, Uint32.class);
    }

    public static Sha3Function load(String contractAddress, Web3j web3j, Credentials credentials, GasProvider contractGasProvider) {
        return new Sha3Function(contractAddress, web3j, credentials, contractGasProvider);
    }

    public static Sha3Function load(String contractAddress, Web3j web3j, TransactionManager transactionManager, GasProvider contractGasProvider) {
        return new Sha3Function(contractAddress, web3j, transactionManager, contractGasProvider);
    }
}
